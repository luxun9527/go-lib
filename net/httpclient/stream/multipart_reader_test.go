package stream

import (
	"io"
	"log"
	"os"

	"net/http"
	"testing"
)

func TestMultiPart1(t *testing.T) {
	rd, err := NewMultipartReaderWriter()
	if err := rd.WriteFiled("key", "value"); err != nil {
		log.Panicf("WriteFiled failed err %v", err)
	}

	fs1, err := os.Open("example.txt")
	if err != nil {
		log.Panicf("err = %v", err)
	}
	defer fs1.Close()
	if err := rd.WriteFileField("example.txt", "file1", fs1); err != nil {
		log.Panicf("WriteFileField1 failed err %v", err)
	}
	fs2, err := os.Open("example.txt")
	if err != nil {
		log.Panicf("err = %v", err)
	}
	defer fs2.Close()

	if err := rd.WriteFileField("example.txt", "file2", fs2); err != nil {
		log.Panicf("WriteFileField2 failed err %v", err)

	}
	_ = rd.Close()
	data, err := io.ReadAll(rd)
	if err != nil {
		log.Panicf("ReadAll failed err %v", err)
	}
	log.Println(string(data))
}

func TestMultiPart2(t *testing.T) {
	rd, err := NewMultipartReaderWriter()
	if err := rd.WriteFiled("key", "value"); err != nil {
		log.Panicf("WriteFiled failed err %v", err)
	}

	fs1, err := os.Open("example.txt")
	if err != nil {
		log.Panicf("err = %v", err)
	}
	defer fs1.Close()
	if err := rd.WriteFileField("example.txt", "file1", fs1); err != nil {
		log.Panicf("WriteFileField1 failed err %v", err)
	}
	fs2, err := os.Open("example.txt")
	if err != nil {
		log.Panicf("err = %v", err)
	}
	defer fs2.Close()

	if err := rd.WriteFileField("example.txt", "file2", fs2); err != nil {
		log.Panicf("WriteFileField2 failed err %v", err)

	}
	_ = rd.Close()

	url := "http://localhost:10011/" // 服务器URL

	req, err := http.NewRequest("POST", url, rd)
	if err != nil {
		log.Panicf("new request error %v", err)
	}
	req.Header.Set("Content-Type", rd.contentType)
	req.TransferEncoding = []string{"chunked"}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("do requeset  error %v", err)
	}
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read error %v", err)
	}
	log.Println(string(d))
}
func TestMultServer(t *testing.T) {
	if err := http.ListenAndServe(":10011", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// 获取文件
		file1, handler, err := request.FormFile("file1")
		if err != nil {
			http.Error(writer, "Error retrieving file", http.StatusInternalServerError)
			return
		}

		defer file1.Close()
		data, _ := io.ReadAll(file1)
		log.Printf("file size: %v content %v filename %v size %v", len(data), string(data), handler.Filename, handler.Size)
		// 获取文件
		file2, handler, err := request.FormFile("file2")
		if err != nil {
			http.Error(writer, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file2.Close()
		data, _ = io.ReadAll(file2)
		log.Printf("file size: %v content %v filename %v size %v", len(data), string(data), handler.Filename, handler.Size)
		value := request.FormValue("key")
		log.Printf("key: %v", value)

		writer.Write([]byte("hello world"))
	})); err != nil {
		log.Panicf("http server error: %v", err)
	}
}
