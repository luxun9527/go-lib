package proxy

import (
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