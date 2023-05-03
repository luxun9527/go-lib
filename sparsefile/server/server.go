package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"os"
)

func main() {
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept failed err %v \n", err)
			continue
		}

		buffer := Buffer{
			buf:  bytes.NewBuffer(make([]byte, 0, 1024*1024*40)),
			path: "/root/testsparse",
		}
		buffer.hande(conn)

	}
}

type Buffer struct {
	buf  *bytes.Buffer
	path string
}

//要考虑一次读不完一条河一次读出多条的情况。
func (c *Buffer) hande(conn net.Conn) {
	buf := make([]byte, 1024*1024*10)
	fd, err := os.OpenFile(c.path, os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return
	}
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read message failed err %v", err)
		}
		c.buf.Write(buf[:n])
		for {
			if c.buf.Len() < 16 {
				continue
			}
			remainBytes := c.buf.Bytes()
			offset := binary.BigEndian.Uint64(remainBytes[:8])
			size := binary.BigEndian.Uint64(remainBytes[8:16])
			if uint64(c.buf.Len()-16) < size {
				continue
			}
			if _, err := fd.WriteAt(remainBytes[16:size], int64(offset)); err != nil {
				return
			}
			//每次执行完做一个拷贝，不然一直引用该切片会让其一直变大
			cache := make([]byte, len(remainBytes)-16-int(size))
			copy(cache, remainBytes[16+size:])
			c.buf = bytes.NewBuffer(cache)
		}
	}

}
