package pool

import (
	"go.uber.org/atomic"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"testing"
	"time"
)

var _httpCli = &http.Client{
	Timeout: time.Duration(15) * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        10,                 //最大空闲
		MaxIdleConnsPerHost: 5,                  //每个host最多保持多少个空闲连接， 如果连接数超过MaxIdleConnsPerHost 则会关闭多余的连接。
		MaxConnsPerHost:     5,                  //MaxConnPerHost 10 决定了每个host最大的连接数，包括正在使用的，正在建立连接的，空闲的，决定了最大并发请求。超过则会阻塞
		IdleConnTimeout:     1000 * time.Second, //空闲的连接超时时间，当超过这个时间则会关闭空闲的连接
	},
	//5个goroutine 并发请求，会有两个并发，其他三个阻塞，
}

//var _httpCli = http.DefaultClient

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
		http.ListenAndServe("0.0.0.0:9999", nil)
	}()
	go func() {
		for {
			for i := 0; i < 50; i++ {
				go get("http://127.0.0.1:9090")
			}
			time.Sleep(time.Second * 3)
		}

	}()

	select {}
}

func TestInitServer(t *testing.T) {
	//https://www.cnblogs.com/paulwhw/p/15972645.html
	//https://www.jianshu.com/p/43bb39d1d221
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()
	var (
		lock          sync.RWMutex
		m             = make(map[string]int, 10)
		receivedCount atomic.Int32
	)

	go func() {
		for {
			time.Sleep(time.Second * 3)
			lock.Lock()
			log.Printf("client size %v request count %v", len(m), receivedCount.Load())
			lock.Unlock()
		}
	}()
	log.Printf("server start %v", 9090)
	if err := http.ListenAndServe(":9090", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("helloworld"))
		receivedCount.Inc()
		lock.Lock()

		value, ok := m[request.RemoteAddr]
		if !ok {
			m[request.RemoteAddr] = 1
		} else {
			value++
			m[request.RemoteAddr] = value
		}
		lock.Unlock()
		select {}

	})); err != nil {
		log.Printf("server start error %v", err)
	}

}
