package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestClient(t *testing.T) {
	//writeChunked()
	readChunked()
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
