package test

import (
	"fmt"
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
	conn, err := net.Dial("tcp", "127.0.0.1:20001")
	if err != nil {
		t.Error(err)
		return
	}
	for {
		time.Sleep(time.Second * 5)
		conn.Write([]byte("abc"))

	}
}
func TestServer1(t *testing.T) {
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
		go func() {
			for {
				buf := make([]byte, 100)
				n, err := conn.Read(buf)
				if err != nil {
					log.Println(err)
				}
				log.Println(string(buf[:n]))
			}
		}()
	}
}
