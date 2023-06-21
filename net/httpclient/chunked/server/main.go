package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {

	/* Net listener */
	n := "tcp"
	addr := ":9094"
	l, err := net.Listen(n, addr)
	if err != nil {
		panic("AAAAH")
	}

	/* HTTP server */
	s := http.Server{
		Handler: http.HandlerFunc(handleRead),
	}
	s.Serve(l)

}

type Message struct {
	Offset string `json:"offset"`
	Data   []byte `json:"data"`
}

func handleRead(w http.ResponseWriter, req *http.Request) {
	buf := make([]byte, 32*1024)
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
func handleWrite(w http.ResponseWriter, req *http.Request) {
	flusher := w.(http.Flusher)

	w.Header().Set("X-Content-Type-Options", "nosniff")
	for i := 1; i <= 20; i++ {
		w.Write([]byte("test111"))
		flusher.Flush() // Trigger "chunked" encoding and send a chunk...
		time.Sleep(1 * time.Second)
	}
	log.Println("finish")

}
