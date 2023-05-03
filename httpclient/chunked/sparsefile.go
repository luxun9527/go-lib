package sparsefile

import (
	"context"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/errors/gerror"
	"github.com/spf13/cast"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Option func(*SparseFile)

func SetProgressHandler(f1 func(progress int64)) Option {
	return func(s *SparseFile) {
		s.ProgressHandler = f1
	}
}

type SparseFile struct {
	file            *os.File
	ctx             context.Context
	lastHole        int64
	Size            int64
	lastProcess     int64
	ProgressHandler func(progress int64)
}

func (s *SparseFile) Close() error {
	return nil
}
func (s *SparseFile) GetSize() int64 {
	blocks := s.GetValidBlock()
	var size int64
	for _, v := range blocks {
		size += v.Size
	}
	return size
}

func NewSparseFile(ctx context.Context, fs *os.File, opts ...Option) *SparseFile {
	fileInfo, _ := fs.Stat()

	sf := &SparseFile{
		file: fs,
		ctx:  ctx,
		Size: fileInfo.Size(),
	}
	for _, opt := range opts {
		opt(sf)
	}
	return sf

}

type ValidBlock struct {
	//偏移
	Offset int64
	//有效块的大小
	Size int64
}

// CopySpareFile 将稀疏文件有效的块拷贝到目的地 使用cp命令更好
func (s *SparseFile) CopySpareFile(distPath string, buf []byte) error {

	vbList := s.GetValidBlock()

	distFs, err := os.OpenFile(distPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}
	defer distFs.Close()
	//元数据区在文件的结尾，记录了文件所在的扇区，如果拷贝有效的block然后顺序存到文件中，相当于改变了文件的结构，文件系统是无法识别的,必须要保持和原来一样的结构。
	for _, v := range vbList {
		if _, err := s.file.Seek(v.Offset, io.SeekStart); err != nil {
			return err
		}
		if _, err := distFs.Seek(v.Offset, io.SeekStart); err != nil {
			return err
		}
		limitReader := io.LimitReader(s, v.Size)
		if err := copyBlock(limitReader, buf, distFs); err != nil {
			return err
		}
	}
	return nil

}

// CopyToRemote 拷贝稀疏文件到远程，思路使用Transfer-Encoding chunked的方式上传，前20个字节表示offset,后面的数据表示真实数据。
func (s *SparseFile) CopyToRemote(url *url.URL, authInfo map[string][]string) error {
	tr := http.DefaultTransport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	tr.WriteBufferSize = 32 * 1024
	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}

	req := &http.Request{
		Method:        "POST",
		URL:           url,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Body:          s,
	}
	req.Header = make(http.Header, len(authInfo))
	for k, v := range authInfo {
		req.Header.Add(k, strings.Join(v, ", "))
	}
	//req.Header.Add("size", cast.ToString(s.GetSize()))
	if _, err := client.Do(req); err != nil {
		return err
	}
	return nil

}

// ReceiveSparseFile 接收
func ReceiveSparseFile(c *gin.Context, distPath string, processHandler func(int64)) error {
	buf := make([]byte, 4*1024)
	distFs, err := os.OpenFile(distPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	size := cast.ToInt64(c.GetHeader("size"))
	defer c.Request.Body.Close()
	var received, lastProcess int64
	//log.Println("size", size)

	reader := httputil.NewChunkedReader(c.Request.Body)
	for {
		n, err := reader.Read(buf)
		log.Printf("client receive %v err %v", n, err)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if processHandler != nil {
			received += int64(n - 20)
			f := cast.ToFloat64(received) / cast.ToFloat64(size) * 100
			//当进度没有改变的时候不处理
			if cast.ToInt64(f) != lastProcess {
				processHandler(cast.ToInt64(f))
				lastProcess = cast.ToInt64(f)
			}
		}
		b1 := buf[:n]
		offset := strings.TrimSpace(string(b1[:20]))
		offset = strings.ReplaceAll(offset, string([]byte{0}), "")
		if len(offset) > 0 && cast.ToInt64(offset) != 0 {
			log.Println("offset", offset)
			if _, err := distFs.Seek(cast.ToInt64(offset), io.SeekStart); err != nil {
				return gerror.Wrap(err, "seek invoke")
			}
		}
		if _, err := distFs.Write(b1[20:]); err != nil {
			return err
		}

	}
	return nil
}

func (s *SparseFile) Read(p []byte) (int, error) {

	start, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	//当进度改变处理
	if s.ProgressHandler != nil {
		f := cast.ToFloat64(start) / cast.ToFloat64(s.Size) * 100
		//当进度没有改变的时候不处理
		if cast.ToInt64(f) != s.lastProcess {
			s.ProgressHandler(cast.ToInt64(f))
			s.lastProcess = cast.ToInt64(f)
		}

	}
	//如果跳到文件的结尾表示结束
	if start == s.Size {
		return 0, io.EOF
	}
	buf := make([]byte, 4076)
	//https://www.zhihu.com/question/407305048
	//SEEK_HOLE的意思就是从offset开始找，找到大于等于offset的第一个Hole开始的地址。如果offset指在一个Hole的中间，那就返回offset。如果offset后面再没有更多的hole了，那就返回文件结尾。
	data, err := s.file.Seek(start, unix.SEEK_DATA)
	if err != nil {
		return 0, err
	}
	//SEEK_DATA的意思很明确，就是从指定的offset开始往后找，找到在大于等于offset的第一个不是Hole的地址。如果offset正好指在一个DATA区域的中间，那就返回offset。
	hole, err := s.file.Seek(data, unix.SEEK_HOLE)
	if err != nil {
		return 0, err
	}

	if _, err := s.file.Seek(data, io.SeekStart); err != nil {
		return 0, err
	}
	size := hole - data
	//如果size 小于我们定义的buf,就将buf修改到到siz的长度。
	if size < int64(len(buf))-20 {
		r := size + 20
		buf = buf[:r]
	}
	//如果hole的位置不等于上个位置则将偏移的位置加上
	if hole != s.lastHole {
		log.Println("offset", data)
		r := cast.ToString(data)
		copy(buf, r)
		s.lastHole = hole
	}

	select {
	case <-s.ctx.Done():
		return 0, s.ctx.Err()
	default:
		n, err := s.file.Read(buf[20:])
		if err != nil {
			return 0, err
		}
		copy(p, buf[:n+20])
		log.Println("send n ", n)
		return n + 20, nil
	}

}

//GetValidBlock 获取有效的块
func (s *SparseFile) GetValidBlock() []*ValidBlock {
	old, _ := s.file.Seek(0, io.SeekCurrent)

	stat, _ := s.file.Stat()
	lastOffset := stat.Size()
	var data, hole int64
	m := make([]*ValidBlock, 0, 10)
	for {
		data, _ = s.file.Seek(hole, unix.SEEK_DATA)
		hole, _ = s.file.Seek(data, unix.SEEK_HOLE)
		vb := &ValidBlock{
			Offset: data,
			Size:   hole - data,
		}
		m = append(m, vb)
		if hole == lastOffset {
			break
		}

	}
	s.file.Seek(old, io.SeekStart)

	return m
}

//拷贝块
func copyBlock(fs io.Reader, buf []byte, dist *os.File) error {
	for {
		n, err := fs.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		_, err = dist.Write(buf[:n])
		if err != nil {
			return err
		}

	}
	return nil
}
