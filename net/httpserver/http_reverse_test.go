package httpserver

import (
	"log"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"net/url"
	"testing"
	"time"
)

// NewProxy creates a new reverse proxy
func NewProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 10,               //每个host最多保持多少个空闲连接， 如果连接数超过MaxIdleConnsPerHost 则会关闭多余的连接。
		MaxConnsPerHost:     10,               //MaxConnPerHost 2 决定了每个host最大的连接数，包括正在使用的，正在建立连接的，空闲的，决定了最大并发请求。超过则会阻塞
		IdleConnTimeout:     90 * time.Second, //空闲的连接超时时间，当超过这个时间则会关闭空闲的连接

	}

	return proxy, nil
}

func TestReverse(t *testing.T) {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()
	target := "http://localhost:8081" // 将请求转发到的目标服务器
	proxy, err := NewProxy(target)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	log.Println("Starting proxy server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
func TestStartServer(t *testing.T) {
	log.Println("Starting target server on :8081")
	if err := http.ListenAndServe(":8081", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("url %v", request.URL.String())
		select {}
	})); err != nil {
		log.Fatal(err)
	}
}
