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
		TLSHandshakeTimeout:    10 * time.Second,
		MaxIdleConns:          1,
		MaxIdleConnsPerHost:    1,                //每个host最多保持多少个空闲连接， 如果连接数超过MaxIdleConnsPerHost 则会关闭多余的连接。
		MaxConnsPerHost:        2,                //MaxConnPerHost 2 决定了每个host最大的连接数，包括正在使用的，正在建立连接的，空闲的，决定了最大并发请求。超过则会阻塞
		IdleConnTimeout:        90 * time.Second, //空闲的连接超时时间，当超过这个时间则会关闭空闲的连接
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