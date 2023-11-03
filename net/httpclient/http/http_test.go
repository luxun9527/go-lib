package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestClientNormal(t *testing.T) {
	// 创建请求对象
	// 创建http客户端
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8888/download", bytes.NewBuffer([]byte(`{"code":"1","value":1"}`)))
	if err != nil {
		fmt.Println("请求对象创建失败：", err)
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
	if _, err := client.Do(req); err != nil {
		log.Println("err",err)
	}
	body, err := io.ReadAll(req.Body)
	if err!=nil{
		log.Println("err",err)
	}
	log.Println(string(body))

	// 请求头
	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "MyGoApp/1.0",
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送HTTP请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 输出响应内容
	fmt.Println("响应状态码:", resp.StatusCode)
	fmt.Println("响应内容:", string(body))
}


func TestServerNormal(t *testing.T) {
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Header)
		body, err := io.ReadAll(r.Body)
		if err!=nil{
			log.Println(err)
		}
		log.Println(string(body))
		log.Println()
	})

	log.Print("Listening on localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
func TestServer(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		for i := 1; i <= 10; i++ {
			fmt.Fprintf(w, "Chunk #%d\n", i)
			flusher.Flush() // Trigger "chunked" encoding and send a chunk...
			time.Sleep(500 * time.Millisecond)
		}
	})

	log.Print("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func TestServer2(t *testing.T) {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		for i := 1; i <= 10; i++ {
			fmt.Fprintf(w, "Chunk #%d\n", i)
			flusher.Flush() // Trigger "chunked" encoding and send a chunk...
			time.Sleep(500 * time.Millisecond)
		}
	})

	log.Print("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
