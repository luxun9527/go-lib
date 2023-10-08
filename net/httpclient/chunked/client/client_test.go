package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestClient(t *testing.T)  {
	//writeChunked()
	readChunked()
}
func writeChunked() {
	tr := http.DefaultTransport
	client := &http.Client{
		Transport: tr,
	}
	source :="test,test,test,"
	r := io.NopCloser(strings.NewReader(source))
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:9094",
			Path:   "/",
		},
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Body:          r,
	}

	fmt.Printf("Doing request\n")
	if _, err := client.Do(req);err!=nil{
		log.Println(err)
	}
}
func readChunked() {
	tr := http.DefaultTransport

	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:9094",
			Path:   "/",
		},
	}
	fmt.Printf("Doing request\n")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("err", err)
		return
	}
	buf := make([]byte, 10)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			break
		}
		log.Println("data", string(buf[:n]))
	}
}
