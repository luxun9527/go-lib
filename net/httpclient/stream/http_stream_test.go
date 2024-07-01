package stream

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"
)

/*
旨在弄清楚http的流相关的概念
包括
1、文件上传
2、chunk流模式
3、go sdk作为客户端服务端文件上传
4、go sdk go客户端服务端文件直接在request body读写数据。
*/

func TestServer(t *testing.T) {

	if err := http.ListenAndServe(":10009", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		buf := make([]byte, 4096)
		for {

			n, err := request.Body.Read(buf)
			if err != nil {
				log.Printf("read error: %v", err)
				break
			}
			log.Printf(string(buf[:n]))
		}
		flusher := writer.(http.Flusher)
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		for i := 1; i <= 20; i++ {
			writer.Write([]byte("test111"))
			flusher.Flush() // Trigger "chunked" encoding and send a chunk...
			time.Sleep(1 * time.Second)
		}

	})); err != nil {
		log.Panicf("http server error: %v", err)
	}
}

func TestClientRequestBody(t *testing.T) {

	req, _ := http.NewRequest(http.MethodPost, "http://localhost:10009", NewLimitedDataReader(200, 10))
	http.DefaultClient.Transport = &http.Transport{
		WriteBufferSize: 1024 * 8,
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("http client request error: %v", err)
	}

	data, err := io.ReadAll(resp.Body)
	log.Printf("resp: %v", string(data))
	defer resp.Body.Close()
}

func TestMultipartServer(t *testing.T) {

	if err := http.ListenAndServe(":10010", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 获取文件
		file, handler, err := request.FormFile("file")
		if err != nil {
			http.Error(writer, "Error retrieving file", http.StatusInternalServerError)
			return
		}

		defer file.Close()
		data, err := io.ReadAll(file)
		log.Printf("file size: %v content %v filename %v size %v", len(data), string(data), handler.Filename, handler.Size)

		writer.Write([]byte("hello world"))
	})); err != nil {
		log.Panicf("http server error: %v", err)
	}
}

func TestMultipartChunked(t *testing.T) {
	filePath := "example.txt"        // 要上传的文件路径
	url := "http://localhost:10010/" // 服务器URL

	// 创建一个缓冲区来存放multipart/form-data内容
	pr, pw := io.Pipe()

	//var requestBody bytes.Buffer
	writer := multipart.NewWriter(pw)
	go func() {
		defer pw.Close()
		defer writer.Close()
		// 添加文件字段
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// 创建文件字段
		part, err := writer.CreateFormFile("file", filePath)
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return
		}

		// 将文件内容复制到 part 中
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println("Error copying file:", err)
			return
		}
		// 关闭multipart writer，以便添加结束边界
		err = writer.Close()
		if err != nil {
			fmt.Println("Error closing writer:", err)
			return
		}
	}()

	// 创建请求
	req, err := http.NewRequest("POST", url, pr)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.TransferEncoding = []string{"chunked"}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(respBody))
}

func TestMultipart(t *testing.T) {
	filePath := "example.txt"        // 要上传的文件路径
	url := "http://localhost:10010/" // 服务器URL

	// 创建一个缓冲区来存放multipart/form-data内容
	var buf = bytes.NewBuffer(make([]byte, 0, 1024))
	//var requestBody bytes.Buffer
	writer := multipart.NewWriter(buf)
	if err := writer.WriteField("key", "value"); err != nil {
		return
	}
	// 添加文件字段
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	_, err = io.Copy(part, file)

	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	// 关闭multipart writer，以便添加结束边界
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(respBody))
}
