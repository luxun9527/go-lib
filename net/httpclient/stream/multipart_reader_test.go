package stream

import (
	"io"
	"log"

	"net/http"
	"testing"
)

func TestMultiPart(t *testing.T) {
	reader, err := NewMultiparkReader()
	if err != nil {
		log.Panicf("NewMultipartReader failed, err:%v\n", err)
	}
	url := "http://localhost:10010/" // 服务器URL

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Panicf("err %v", err)
	}
	req.Header.Set("Content-Type", reader.contentType)

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
