package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestProxy(t *testing.T)  {
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
}

func TestExample(t *testing.T) {
	// Create a Resty Client
	client := resty.New().SetRetryCount(4)

	resp, err := client.R().
		EnableTrace().
		Get("https://httpbin.org/get")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	fmt.Println("Request Trace Info:")
	ti := resp.Request.TraceInfo()
	fmt.Println("  DNSLookup     :", ti.DNSLookup)
	fmt.Println("  ConnTime      :", ti.ConnTime)
	fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	fmt.Println("  ServerTime    :", ti.ServerTime)
	fmt.Println("  ResponseTime  :", ti.ResponseTime)
	fmt.Println("  TotalTime     :", ti.TotalTime)
	fmt.Println("  IsConnReused  :", ti.IsConnReused)
	fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
	fmt.Println("  RequestAttempt:", ti.RequestAttempt)
	fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
}

