package main

import (
	"errors"
	"flag"
	"go-lib/net/httpclient/stream"
	"go.uber.org/atomic"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	_url = "http://127.0.0.1:10010"
)

var (
	mode               = flag.String("mode", "client", "客户端还是服务端")
	c                  = flag.Int64("c", 70, "并发数")
	transferBytesCount = flag.Int64("t", 1024*1024*50, "每个并发数传输的字节数")
	targetUrl          = flag.String("targetUrl", "http://47.113.223.16:9993", "url")
)

func main() {
	flag.Parse()
	if *mode == "client" {
		client()
	} else {
		server()
	}
}

func server() {
	var (
		receivedBytes           atomic.Int64
		receiveReqCount         atomic.Int64
		lastSecondReceivedBytes atomic.Int64
	)
	go func() {
		for {
			time.Sleep(time.Second)
			c := receivedBytes.Load()
			perSec := c - lastSecondReceivedBytes.Load()
			log.Printf("receivedBytes:%d, receiveReqCount:%d perSecond:%d", c, receiveReqCount.Load(), perSec)
			lastSecondReceivedBytes.Store(c)
		}
	}()
	if err := http.ListenAndServe(":9993", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		receiveReqCount.Inc()
		buf := make([]byte, 1024*500)
		for {
			n, err := request.Body.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				log.Printf("read body error: %v", err)
			}
			receivedBytes.Add(int64(n))
		}

	})); err != nil {
		log.Panicf("http server error: %v", err)
	}

}

var (
	requestCount            atomic.Int64
	failedCount             atomic.Int64
	transferBytes           atomic.Int64
	lastSecondTransferBytes atomic.Int64
	debug                   bool
)

type ReaderStat struct {
	r *stream.LimitedDataReader
}

func (d *ReaderStat) Read(p []byte) (int, error) {
	n, err := d.r.Read(p)
	transferBytes.Add(int64(n))
	return n, err
}
func newReaderStat(size int64, perReadSize int64) *ReaderStat {
	return &ReaderStat{
		r: stream.NewLimitedDataReader(size, perReadSize),
	}
}

func client() {

	go func() {
		for {
			time.Sleep(time.Second)
			c := transferBytes.Load()
			perSecondCount := c - lastSecondTransferBytes.Load()
			log.Printf("request count: %d, failed count: %d, transfer bytes: %d kb, transfer per second %d kb", requestCount.Load(), failedCount.Load(), transferBytes.Load()/1024, perSecondCount/1024)
			lastSecondTransferBytes.Store(c)
		}
	}()
	var group sync.WaitGroup
	cli := http.DefaultClient
	for i := int64(0); i < *c; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			reader := newReaderStat(*transferBytesCount, 1024*10)
			_, err := cli.Post(*targetUrl, http.MethodPost, reader)
			requestCount.Inc()
			if err != nil {
				log.Printf("request error: %v", err)
				failedCount.Inc()

			}
		}()

	}
	group.Wait()
}
