package server

import (
	"io"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T)  {

	n := "tcp"
	addr := ":9094"
	l, err := net.Listen(n, addr)
	if err != nil {
		panic(err)
	}

	s := http.Server{
		//Handler: http.HandlerFunc(write),
		Handler: http.HandlerFunc(read),
	}
	s.Serve(l)

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
		w.Write([]byte("test111"))
		flusher.Flush() // Trigger "chunked" encoding and send a chunk...
		time.Sleep(1 * time.Second)
	}
	log.Println("finish")

}
