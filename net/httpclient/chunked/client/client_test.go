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

func TestClient(t *testing.T) {
	writeChunked()
	//readChunked()
}

func readChunked() {

}
func writeChunked1() {
	tr := http.DefaultTransport
	client := &http.Client{
		Transport: tr,
	}
	source := "test,test,test,"
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
	if _, err := client.Do(req); err != nil {
		log.Println(err)
	}
}

func writeChunked() {
	tr := http.DefaultTransport

	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	b := strings.NewReader("test")
	cb := io.NopCloser(b)
	req := &http.Request{
		Method: "PUT",
		URL: &url.URL{
			Scheme: "http",
			Host:   "192.168.2.159:9001",
			Path:   "/test/test?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20240612T153331Z&X-Amz-SignedHeaders=host&X-Amz-Expires=3599&X-Amz-Credential=admin%2F20240612%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=43bc225b6baa919b1ba76464e838be25296a1e46b66847e87dd8ff07615ba569",
		},
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Body:          cb,
		//Header:        map[string][]string{"Content-Length": []string{"111"}},
	}

	fmt.Printf("Doing request\n")
	resp, err := client.Do(req)
	fmt.Printf("Done request. Err: %v\n", err)
	data, err := io.ReadAll(resp.Body)
	log.Printf("%v", string(data))
}
