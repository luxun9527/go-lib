package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close() // 关闭连接
	for {
		buf := make([]byte, 128)

		n, err := conn.Read(buf) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", recvStr)
	}
}
func TestServer(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:20001")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}

	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		//conn.SetWriteDeadline()
		go process(conn) // 启动一个goroutine处理连接
	}
}
func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.2.99:8080")
	if err != nil {
		t.Error(err)
		return
	}
	for {
		conn.Write([]byte("abc"))
		time.Sleep(time.Second * 20)
	}

}
func TestClient1(t *testing.T) {
	for i := 0; i < 10000; i++ {
		go func() {
			conn, err := net.Dial("tcp", "47.113.223.16:9993")
			if err != nil {
				t.Error(err)
				return
			}
			for {
				time.Sleep(time.Second * 5)
				conn.Write([]byte("abc"))

			}
		}()
	}
	select {}

}
func TestServer1(t *testing.T) {
	listen, err := net.Listen("tcp", "0.0.0.0:20001")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}

	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go func() {
			buffer := bytes.NewBuffer(make([]byte, 0, 1024))
			if _, err := io.Copy(buffer, conn); err != nil {
				log.Printf("io.Copy err: %v", err)
			}
			log.Println(buffer.String())
		}()
	}
}

// 测试当客户端发送，不进行操作服务端是否会断开连接
func TestTcpCli(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:9094")
	if err != nil {
		t.Error(err)
		return
	}
	data := `GET /test HTTP/1.1
Language: zh-CN
gexToken: undefined
User-Agent: Apifox/1.0.0 (https://apifox.com)
Accept: */*
Host: localhost:20001
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
`
	_, err = conn.Write([]byte(data))
	if err != nil {
		log.Printf("conn.Write err: %v", err)
	}
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	if _, err := io.Copy(buffer, conn); err != nil {
		log.Printf("io.Copy err: %v", err)
	}
	log.Printf("data %+v", buffer.String())

	time.Sleep(time.Second * 100)
}
