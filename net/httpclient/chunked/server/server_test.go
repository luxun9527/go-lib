package server

import (
	"fmt"
	"io"
	"log"

	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/read", read)
	mux.HandleFunc("/write", write)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicf("ListenAndServe: %v", err)
	}
}

func read(w http.ResponseWriter, req *http.Request) {
	buf := make([]byte, 10)
	for {
		n, err := req.Body.Read(buf)
		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			break
		}
		//time.Sleep(time.Second)
		b1 := buf[:n]
		log.Println("data", string(b1))
	}
	log.Println("finish")

}
func write(w http.ResponseWriter, req *http.Request) {
	flusher := w.(http.Flusher)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	for i := 1; i <= 20; i++ {
		w.Write([]byte(fmt.Sprintf("hello world %d\n", i)))
		flusher.Flush()
		time.Sleep(1 * time.Second)
	}
	log.Println("finish")

}
