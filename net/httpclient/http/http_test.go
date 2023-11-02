package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	params := []byte(`{"key1": "value1", "key2": "value2"}`)
	// 请求头
	headers := map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "MyGoApp/1.0",
	}
	// 创建HTTP请求
	req, err := http.NewRequest("POST", "http://127.0.0.1:8888/api", bytes.NewBuffer(params))
	if err != nil {
		fmt.Println("创建HTTP请求失败:", err)
		return
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 发送HTTP请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送HTTP请求失败:", err)
		return
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应内容失败:", err)
		return
	}
	// 输出响应内容
	fmt.Println("响应状态码:", resp.StatusCode)
	fmt.Println("响应内容:", string(body))
}
func TestServer11(t *testing.T) {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {

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
