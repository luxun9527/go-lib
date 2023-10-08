



# http客户端实践

## http代理

```go
	// 创建一个Clash代理的HTTP Transport
	proxyURL, err := url.Parse("http://192.168.11.85:7890") // 替换为你的Clash代理地址和端口
	if err != nil {
		log.Fatal(err)
	}
	// 创建一个Transport，使用Clash代理
	transport := &http.Transport{
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     50,
		IdleConnTimeout:     30 * time.Second,
		Proxy: http.ProxyURL(proxyURL),
	}
	//resty.New().
	response, err := resty.New().SetTransport(transport).R().Get("https://is.gd/create.php?format=json&url=www.baidu.com")
	if err!=nil{
		log.Fatal(err)
	}
	body := response.Body()
	log.Println(string(body))
```



## http连接池配置

```go
package pool

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// TODO 计算客户端host(IP+Port)的数量
var m = make(map[string]int)

var ch = make(chan string, 10)

// TODO 计算链接数量
func count() {
	for s := range ch {
		m[s]++
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	logrus.Info(r.RemoteAddr) // TODO 最后打印的是 remoteAddr
	ch <- r.RemoteAddr
	// time.Sleep(time.Second)
	w.Write([]byte("helloworld"))
}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func graceClose() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(ch)
		time.Sleep(time.Second)
		spew.Dump(m)
		os.Exit(0)
	}()
}

var _httpCli = &http.Client{
	Timeout: time.Duration(15) * time.Second,

	Transport: &http.Transport{
		//MaxIdleConns:          1,
		MaxIdleConnsPerHost:   1, //每个host最多保持多少个空闲连接， 如果连接数超过MaxIdleConnsPerHost 则会关闭多余的连接。
		MaxConnsPerHost:       2,//MaxConnPerHost 2 决定了每个host最大的连接数，包括正在使用的，正在建立连接的，空闲的，决定了最大并发请求。超过则会阻塞
		IdleConnTimeout:       90 * time.Second, //空闲的连接超时时间，当超过这个时间则会关闭空闲的连接
		TLSHandshakeTimeout:   10 * time.Second,
		//ExpectContinueTimeout: 1 * time.Second,
	},
	//5个goroutine 并发请求，会有两个并发，其他三个阻塞，
}

func get(url string) {
	resp, err := _httpCli.Get(url)
	if err != nil {
		// do nothing
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// do nothing
		return
	}
}


func TestLong(t *testing.T) {
	go func() {
		for i := 0; i < 50; i++ {
			go get("http://127.0.0.1:8087")
		}
		time.Sleep(time.Second*3)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond*500)
			go get("http://127.0.0.1:8087")
		}
	}()

	time.Sleep(time.Second * 100)

}

func TestInitServer(t *testing.T){
	//https://www.cnblogs.com/paulwhw/p/15972645.html
	//https://www.jianshu.com/p/43bb39d1d221
	graceClose()
	go count()
	port := flag.Int("port", 8087, "")
	flag.Parse()

	logrus.Println("Listen port:", *port)

	http.HandleFunc("/", home)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic(err)
	}
}
```

**核心配置**

1. ​	MaxIdleConnsPerHost:   1, //每个host最多保持多少个空闲连接， 如果连接数超过MaxIdleConnsPerHost 则会关闭多余的连接。
2. ​	MaxConnsPerHost:       2,//MaxConnPerHost 2 决定了每个host最大的连接数，包括正在使用的，正在建立连接的，空闲的，决定了最大并发请求。超过则会阻塞
3. ​	IdleConnTimeout:       90 * time.Second, //空闲的连接超时时间，当超过这个时间则会关闭空闲的连接



```go
func TestLong(t *testing.T) {
	go func() {
		for i := 0; i < 50; i++ {
			go get("http://127.0.0.1:8087")
		}
		time.Sleep(time.Second*3)
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond*500)
			go get("http://127.0.0.1:8087")
		}
	}()

	time.Sleep(time.Second * 100)
}
//前50个复用两个连接，请求完成后，sleep 3s 由于设置MaxIdleConnsPerHost：1 空闲的连接为1 则会关闭一个连接，后面的请求会使用另一个连接。
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54202"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54202"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54202"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54202"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54202"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54202"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:05+08:00" level=info msg="127.0.0.1:54202"
time="2023-10-08T10:57:08+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:09+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:09+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:10+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:10+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:11+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:11+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:12+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:12+08:00" level=info msg="127.0.0.1:54201"
time="2023-10-08T10:57:13+08:00" level=info msg="127.0.0.1:54201"
```



## chunked流模式

https://gist.github.com/ZenGround0/af448f56882c16aaf10f39db86b4991e
https://en.wikipedia.org/wiki/Chunked_transfer_encoding
https://juejin.cn/post/6997215152533667876
https://juejin.cn/post/6844903937825488909
https://github.com/golang/go/issues/18407 

```go
package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestClient(t *testing.T)  {
	//writeChunked()
	readChunked()
}
func writeChunked() {
	tr := http.DefaultTransport
	client := &http.Client{
		Transport: tr,
	}

	r := io.NopCloser(strings.NewReader("test,test,test"))

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:9094",
			Path:   "/",
		},
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Body:          r,
	}

	fmt.Printf("Doing request\n")
	if _, err := client.Do(req);err!=nil{
		log.Println(err)
	}
}
func readChunked() {
	tr := http.DefaultTransport

	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:9094",
			Path:   "/",
		},
	}
	fmt.Printf("Doing request\n")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("err", err)
		return
	}
	buf := make([]byte, 10)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			break
		}
		log.Println("data", string(buf[:n]))
	}
}

```

```go
package server

import (
	"io"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T)  {

	n := "tcp"
	addr := ":9094"
	l, err := net.Listen(n, addr)
	if err != nil {
		panic(err)
	}

	s := http.Server{
		Handler: http.HandlerFunc(write),
		//Handler: http.HandlerFunc(read),
	}
	s.Serve(l)

}


func read(w http.ResponseWriter, req *http.Request) {
	buf := make([]byte, 10)
	for {
		n, err := req.Body.Read(buf)
		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			break
		}
		//time.Sleep(time.Second)
		b1 := buf[:n]
		log.Println("data", string(b1))
	}
	log.Println("finish")

}
func write(w http.ResponseWriter, req *http.Request) {
	flusher := w.(http.Flusher)

	w.Header().Set("X-Content-Type-Options", "nosniff")
	for i := 1; i <= 20; i++ {
		w.Write([]byte("test111"))
		flusher.Flush() // Trigger "chunked" encoding and send a chunk...
		time.Sleep(1 * time.Second)
	}
	log.Println("finish")

}

```

