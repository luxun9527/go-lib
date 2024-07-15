package httpserver

import (
	"bufio"
	"fmt"
	"go.uber.org/atomic"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	n := "tcp"
	addr := ":9094"
	l, err := net.Listen(n, addr)
	if err != nil {
		log.Panicf("Error listening: %v", err.Error())
	}

	s := http.Server{
		Addr:                         "",
		Handler:                      http.HandlerFunc(handle),
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  time.Second * 3,
	}
	s.SetKeepAlivesEnabled(false)
	if err := s.Serve(l); err != nil {
		log.Panicf("Error serving: %v", err.Error())
	}

}

func handle(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second * 1)
	w.Write([]byte("hello"))
	log.Println("finish")

}

func TestHttpServer2(t *testing.T) {
	if err := http.ListenAndServe("0.0.0.0:10001", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		time.Sleep(time.Second)

		writer.Write([]byte(`{"code":200}`))

	})); err != nil {
		log.Panicf("Error listening: %v", err.Error())
	}
}

func TestHttpCli3(t *testing.T) {
	go func() {
		http.ListenAndServe("0.0.0.0:10003", nil)
	}()
	listener, err := net.Listen("tcp", ":10002")
	if err != nil {
		log.Panicf("Error listening: %v", err.Error())
	}
	var count atomic.Int64
	go func() {
		for {
			time.Sleep(time.Second)
			log.Printf("count: %d", count.Load())
		}
	}()
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			MaxConnsPerHost:     100,
			IdleConnTimeout:     time.Second * 120,
		},
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting: %v", err.Error())
			continue
		} // 创建一个 HTTP 客户端

		go func(conn net.Conn) {

			for {

				// 创建一个 bufio.Reader 用于读取 TCP 连接数据
				reader := bufio.NewReader(conn)
				req, err := http.ReadRequest(reader)
				if err != nil {
					if err != io.EOF {
						log.Println("Error reading request:", err)
					}
					break
				}
				// 修改请求 URL 为目标 URL
				req.URL.Scheme = "http"
				req.URL.Host = "localhost:10001"

				// 删除代理相关的头部
				req.RequestURI = ""
				req.Header.Del("Proxy-Connection")
				// 转发请求到目标 HTTP 服务器
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Error forwarding request:", err)
					return
				}
				defer resp.Body.Close()
				if err := resp.Write(conn); err != nil {
					count.Inc()
					return
				}

			}

		}(conn)
	}
}
func TestHttpCli(t *testing.T) {
	var count atomic.Int64
	for i := 0; i < 50; i++ {
		go func() {
			for i := 0; i < 100; i++ {
				resp, err := http.Get("http://localhost:10002")
				if err != nil {
					log.Printf("Error: %v", err)
				}
				_, err = io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error: %v", err)
				}
				if resp.StatusCode == http.StatusOK {
					count.Inc()
				}
				//log.Printf("data=%v", string(data))
			}
		}()
	}
	go func() {
		for {
			time.Sleep(time.Second)
			log.Println("count=", count.Load())
		}
	}()
	time.Sleep(time.Hour)
}
