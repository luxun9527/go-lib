package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	fs, err := os.OpenFile("test.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("err", err)
		return
	}
	go func() {
		for {
			time.Sleep(time.Second)
			fs.WriteString("aaaaa\n")
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second)
			if err := fs.Truncate(0); err != nil {
				log.Println("err", err)
			}
		}
	}()
	time.Sleep(time.Hour)
}
func f() {
	fs, err := os.OpenFile("/Users/demg/personProject/go-lib/io/file/test.txt", os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	buf := make([]byte, 4096)
	for {
		n, err := fs.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read err", err)
			return
		}
		if err == io.EOF || n == 0 {
			fmt.Println("end")
			break
		}
		fmt.Println("data = ", string(buf[:n]))

	}
}
