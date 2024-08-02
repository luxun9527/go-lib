package main

import (
	"errors"
	"fmt"
	"github.com/gookit/goutil/mathutil"
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
	conn, err := net.DialTimeout("tcp", "192.168.2.109:20001", time.Second*3)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		if _, err := conn.Write([]byte("abc")); err != nil {
			log.Printf("write failed, err:%v\n", err)
		}
		time.Sleep(time.Second * 1)
	}

}
func TestClientLogStash(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.2.159:20010")
	if err != nil {
		t.Error(err)
		return
	}
	for {
		if _, err := conn.Write([]byte(fmt.Sprintf(`{"code":%d}\n`, mathutil.RandInt(100, 1000)))); err != nil {
			log.Printf("write to logstash failed, err:%v\n", err)
		}
		log.Println("send to logstash success")
		time.Sleep(time.Second * 3)
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
			d := make([]byte, 1024)
			for {
				n, err := conn.Read(d)
				if err != nil {
					if errors.Is(err, io.EOF) {
						log.Printf("client close")
						break
					}
					log.Println("read err:", err)
					break
				}
				log.Println(string(d[:n]))
			}

		}()
	}
}
