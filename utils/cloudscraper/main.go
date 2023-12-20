package main

import (
	"bytes"
	"encoding/json"
	hawk "github.com/juiced-aio/hawk-go"
	http "github.com/useflyent/fhttp"
	"github.com/useflyent/fhttp/cookiejar"
	"io"
	"log"
)

// refer Port of HawkAPI's cloudscraper
func main() {
	//https://github.com/juiced-aio/hawk-go
	//url := "https://goo.su/api/links/create"
	//m := map[string]string{"x-goo-api-token": "W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1"}
	//
	b := map[string]interface{}{"alias": "", "is_public": true, "group_id": 1, "url": "https://www.api.com"}
	d, _ := json.Marshal(b)
	//
	//	// 创建http客户端
	//	client := &http.Client{}
	//
	//	// 创建请求对象
	//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	//	if err != nil {
	//		fmt.Println("请求对象创建失败：", err)
	//		return
	//	}
	//
	//	//req.Header.Add("Accept-Encoding", "gzip, deflate")
	//	req.Header.Add("content-type", " application/json")
	//	req.Header.Add("user-agent", "Apifox/1.0.0 (https://apifox.com)")
	//	req.Header.Add("content-length", "71")
	//	req.Header.Add("accept", " */*")
	//	// 设置请求头
	////	req.Header.Set("x-goo-api-token","W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1")
	//	req.Header.Set("x-goo-api-token","W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1")
	// 发送请求
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("请求发送失败：", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//bo, err := io.ReadAll(resp.Body)
	//if err!=nil{
	//	log.Println("err",err)
	//}
	//// 读取响应内容
	//var response =map[string]interface{}{}
	//if err := json.Unmarshal(bo,&response); err != nil {
	//
	//	fmt.Println("响应解析失败：", err,string(bo))
	//	return
	//}

	// Client has to be from fhttp and up to CloudFlare's standards, this can include ja3 fingerprint/http2 settings.
	client := http.Client{}
	// Client also will need a cookie jar.
	cookieJar, _ := cookiejar.New(nil)
	client.Jar = cookieJar
	scraper := hawk.CFInit(client, "YOUR_KEY_HERE", true)

	// You will have to create your own function if you want to solve captchas.
	scraper.CaptchaFunction = func(originalURL string, siteKey string) (string, error) {
		// CaptchaFunction should return the token as a string.
		return "", nil
	}

	req, _ := http.NewRequest("POST", "https://goo.su/api/links/create", bytes.NewBuffer(d))

	req.Header = http.Header{
		"sec-ch-ua":                 {`"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`},
		"sec-ch-ua-mobile":          {`?0`},
		"upgrade-insecure-requests": {`1`},
		"user-agent":                {`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36`},
		"accept":                    {`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`},
		"sec-fetch-site":            {`none`},
		"sec-fetch-mode":            {`navigate`},
		"sec-fetch-user":            {`?1`},
		"sec-fetch-dest":            {`document`},
		"accept-encoding":           {`gzip, deflate`},
		"accept-language":           {`en-US,en;q=0.9`},
		"x-goo-api-token":           {"W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1"},
		http.HeaderOrderKey:         {"sec-ch-ua", "sec-ch-ua-mobile", "upgrade-insecure-requests", "user-agent", "accept", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "accept-encoding", "accept-language"},
		http.PHeaderOrderKey:        {":method", ":authority", ":scheme", ":path"},
	}

	resp, err := scraper.Do(req)
	if err != nil {
		log.Printf("err err=%v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}
func test1() {
	//url := "https://goo.su/api/links/create"
	//m := map[string]string{"x-goo-api-token": "W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1"}
	//
	b := map[string]interface{}{"alias": "", "is_public": true, "group_id": 1,      "url":"https://www.api.com"}
	d, _ := json.Marshal(b)
	//
	//      // 创建http客户端
	//      client := &http.Client{}
	//
	//      // 创建请求对象
	//      req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	//      if err != nil {
	//              fmt.Println("请求对象创建失败：", err)
	//              return
	//      }
	//
	//      //req.Header.Add("Accept-Encoding", "gzip, deflate")
	//      req.Header.Add("content-type", " application/json")
	//      req.Header.Add("user-agent", "Apifox/1.0.0 (https://apifox.com)")
	//      req.Header.Add("content-length", "71")
	//      req.Header.Add("accept", " */*")
	//      // 设置请求头
	////    req.Header.Set("x-goo-api-token","W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1")
	//      req.Header.Set("x-goo-api-token","W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1")
	// 发送请求
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//      fmt.Println("请求发送失败：", err)
	//      return
	//}
	//defer resp.Body.Close()
	//
	//bo, err := io.ReadAll(resp.Body)
	//if err!=nil{
	//      log.Println("err",err)
	//}
	//// 读取响应内容
	//var response =map[string]interface{}{}
	//if err := json.Unmarshal(bo,&response); err != nil {
	//
	//      fmt.Println("响应解析失败：", err,string(bo))
	//      return
	//}

	// Client has to be from fhttp and up to CloudFlare's standards, this can include ja3 fingerprint/http2 settings.
	client :=http.Client{}
	// Client also will need a cookie jar.
	cookieJar, _:= cookiejar.New(nil)
	client.Jar = cookieJar
	scraper:= hawk.CFInit(client, "YOUR_KEY_HERE", true)

	// You will have to create your own function if you want to solve captchas.
	scraper.CaptchaFunction = func(originalURL string, siteKey string) (string, error) {
		// CaptchaFunction should return the token as a string.
		return "", nil
	}

	req, _ := http.NewRequest("POST", "https://goo.su/api/links/create",  bytes.NewBuffer(d))

	req.Header = http.Header{
		"sec-ch-ua":                 {`"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`},
		"sec-ch-ua-mobile":          {`?0`},
		"upgrade-insecure-requests": {`1`},
		"user-agent":                {`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36`},
		"accept":                    {`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`},
		"sec-fetch-site":            {`none`},
		"sec-fetch-mode":            {`navigate`},
		"sec-fetch-user":            {`?1`},
		"sec-fetch-dest":            {`document`},
		"accept-encoding":           {`gzip, deflate`},
		"accept-language":           {`en-US,en;q=0.9`},
		"x-goo-api-token":{"W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1"},
		http.HeaderOrderKey:         {"sec-ch-ua", "sec-ch-ua-mobile", "upgrade-insecure-requests", "user-agent", "accept", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "accept-encoding", "accept-language"},
		http.PHeaderOrderKey:        {":method", ":authority", ":scheme", ":path"},
	}

	resp, err := scraper.Do(req)
	if err!=nil{
		log.Printf("err err=%v",err)
	}
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}