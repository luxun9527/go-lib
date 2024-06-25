package httpserver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
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
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		//IdleTimeout:                  time.Second * 3,
		MaxHeaderBytes: 0,
		TLSNextProto:   nil,
		ConnState:      nil,
		ErrorLog:       nil,
		BaseContext:    nil,
		ConnContext:    nil,
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
		writer.Write([]byte("hello"))
		writer.Write([]byte("target"))
	})); err != nil {
		log.Panicf("Error listening: %v", err.Error())
	}
}

func TestHttpCli3(t *testing.T) {
	listener, err := net.Listen("tcp", ":10002")
	if err != nil {
		log.Panicf("Error listening: %v", err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting: %v", err.Error())
			continue
		}
		go func(conn net.Conn) {
			// 创建一个 HTTP 客户端
			client := &http.Client{}
			// 创建一个 bufio.Reader 用于读取 TCP 连接数据
			reader := bufio.NewReader(conn)
			req, err := http.ReadRequest(reader)
			for {
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
				respBytes, err := httputil.DumpResponse(resp, true)
				if err != nil {
					fmt.Println("Error dumping response:", err)
					return
				}
				_, err = conn.Write(respBytes)
				if err != nil {
					fmt.Println("Error writing response to connection:", err)
					return
				}

				//resp.Body.Close()
			}

		}(conn)
	}
}
