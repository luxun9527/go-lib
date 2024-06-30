package fileupload

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

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
func TestMultipartClient(t *testing.T) {
	filePath := "dist.tar.gz"        // 要上传的文件路径
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
