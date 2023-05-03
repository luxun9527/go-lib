package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"flag"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	s := flag.String("path", "", "文件路径")
	addr := flag.String("addr", "", "目的地的ip和端口")
	flag.Parse()
	if *s == "" || *addr == "" {
		log.Fatalln("path and addr must have a value")
	}
	fd, err := os.Open(*s)
	if err != nil {
		log.Fatalln("path valid", err)
	}
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalln("addr valid", err)
	}
	target := RemoteTarget{conn: conn}
	if err := Copy(context.Background(), fd, target); err != nil {
		log.Printf("copy to remote failed %v", err)
	}
}

type RemoteTarget struct {
	conn net.Conn
}

func (rt RemoteTarget) WriteAt(p []byte, off int64) (n int, err error) {

	buf := bytes.NewBuffer(make([]byte, 0, 16))
	//偏移
	if err := binary.Write(buf, binary.BigEndian, uint64(off)); err != nil {
		return 0, err
	}
	//长度
	if err := binary.Write(buf, binary.BigEndian, uint64(len(p))); err != nil {
		return 0, err
	}
	if _, err := rt.conn.Write(buf.Bytes()); err != nil {
		return 0, err
	}
	//数据
	if _, err := rt.conn.Write(p); err != nil {
		return 0, err
	}
	return 0, nil
}

// Copy  将稀疏文件有效的块拷贝到目的地
func Copy(ctx context.Context, srcFs *os.File, writer io.WriterAt) error {

	curOffset, hole := int64(0), int64(0)
	curHole, lastHole := int64(0), int64(0)
	stat, _ := srcFs.Stat()
	end := stat.Size()

	for {
		buf := make([]byte, 1024*512)
		//如果跳到文件的结尾表示结束
		if curOffset == end {
			return nil
		}

		//https://www.zhihu.com/question/407305048
		//SEEK_DATA的意思很明确，就是从指定的offset开始往后找，找到在大于等于offset的第一个不是Hole的地址。如果offset正好指在一个DATA区域的中间，那就返回offset。
		//不要去处理这个错误，当文件为空或一些异常情况这个地方会报错
		data, _ := srcFs.Seek(curOffset, unix.SEEK_DATA)
		//有时出现hole不是结尾，当data变成0的时候,data会小于上个hole的位置。
		if data < lastHole {
			return nil
		}
		//SEEK_HOLE的意思就是从offset开始找，找到大于等于offset的第一个Hole开始的地址。如果offset指在一个Hole的中间，那就返回offset。如果offset后面再没有更多的hole了，那就返回文件结尾。
		hole, _ = srcFs.Seek(data, unix.SEEK_HOLE)
		//空文件直接返回
		if hole == 0 && data == 0 {
			return nil
		}
		if hole != curHole {
			lastHole = curHole
			curHole = hole
		}
		//跳到数据的区的位置
		curOffset, _ = srcFs.Seek(data, io.SeekStart)

		dataZoneSize := hole - data
		//如果dataZoneSize 小于我们定义的buf,就将buf修改到到dataZoneSize的长度。
		if dataZoneSize < int64(len(buf)) {
			buf = buf[:dataZoneSize]
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := srcFs.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if err == io.EOF {
				return nil
			}
			if _, err := writer.WriteAt(buf[:n], curOffset); err != nil {
				return err
			}
			curOffset += int64(n)
		}
	}

}
