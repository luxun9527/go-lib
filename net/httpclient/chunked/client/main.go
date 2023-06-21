package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

type RWC struct {
	data     []byte
	file     *os.File
	ctx      context.Context
	dataChan chan []byte
	err      chan error
}
type Message struct {
	Offset string `json:"offset"`
	Data   []byte `json:"data"`
}

// 使用append的方案？？
//每次都是读32 * 1024 个字节 我们这边可以固定这个值。
func (r *RWC) Read(p []byte) (int, error) {

	//使用append的方式避免Marshal和UnMarshal大包的cpu消耗
	//定义前20个字节表示偏移的位置
	buf := make([]byte, 10)
	n, err := r.file.Read(buf)
	if err != nil {
		return 0, err
	}
	copy(p, buf[:n])
	return n, nil
}

func (r *RWC) Write() (n int, err error) {

	return 0, err
}
func (r *RWC) Close() error {
	return nil
}

func main() {
	writeChunked()
	//readChunked()
}
func writeChunked() {
	tr := http.DefaultTransport

	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	fs, err := os.Open("S:\\go-lib\\httpclient\\chunked\\client\\test.txt")
	if err != nil {
		log.Println("err", err)
		return
	}

	rwc := RWC{file: fs}
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost:9094",
			Path:   "/",
		},
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Body:          &rwc,
	}

	fmt.Printf("Doing request\n")
	_, err = client.Do(req)
	fmt.Printf("Done request. Err: %v\n", err)
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
func BtrfsSend(path string, urlPath string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := fmt.Sprintf(`btrfs send '%s'`, path)
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", c)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	u, err := url.Parse(urlPath)
	if err != nil {
		return err
	}
	req := &http.Request{
		Method:        "POST",
		URL:           u,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Body:          stdoutPipe,
	}
	tr := http.DefaultTransport

	client := &http.Client{
		Transport: tr,
		Timeout:   0,
	}
	_, err = client.Do(req)
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
