package utils

import (
	"context"
	"golang.org/x/sys/unix"
	"io"
	"os"
)

// Copy  将稀疏文件有效的块拷贝到目的地
func Copy(ctx context.Context, srcPath string, writer io.WriterAt) error {

	srcFs, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFs.Close()
	stat, _ := srcFs.Stat()
	start, end := int64(0), stat.Size()
	buf := make([]byte, 1024*512)
	for {
		//如果跳到文件的结尾表示结束
		if start == end {
			return nil
		}
		//https://www.zhihu.com/question/407305048
		//SEEK_DATA的意思很明确，就是从指定的offset开始往后找，找到在大于等于offset的第一个不是Hole的地址。如果offset正好指在一个DATA区域的中间，那就返回offset。
		data, err := srcFs.Seek(start, unix.SEEK_DATA)
		if err != nil {
			return err
		}
		//SEEK_HOLE的意思就是从offset开始找，找到大于等于offset的第一个Hole开始的地址。如果offset指在一个Hole的中间，那就返回offset。如果offset后面再没有更多的hole了，那就返回文件结尾。
		hole, err := srcFs.Seek(data, unix.SEEK_HOLE)
		if err != nil {
			return err
		}
		//跳到数据的区的位置
		if start, err = srcFs.Seek(data, io.SeekStart); err != nil {
			return err
		}
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
			if _, err := writer.WriteAt(buf[:n], start); err != nil {
				return err
			}
			start += int64(n)
		}
	}

}
