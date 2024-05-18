package error

import (
	"io"
	"log"
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	go InitApiSrv()
	go InitRpcSrv()
	select {}
}
func TestHttpCli(t *testing.T) {
	resp, err := http.Get("http://127.0.0.1:9999/error")
	if err != nil {
		log.Panicf("http get error: %v", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("read resp body error: %v", err)
	}
	log.Printf("resp body: %s", data)

}
