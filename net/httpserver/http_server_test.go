package httpserver

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()

	n := "tcp"
	addr := ":9094"
	l, err := net.Listen(n, addr)
	if err != nil {
		log.Panicf("Error listening: %v", err.Error())
	}

	s := http.Server{
		Addr:                         "",
		Handler:                      http.HandlerFunc(handle),
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  time.Second * 3,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		//IdleTimeout:                  time.Second * 3,
		MaxHeaderBytes: 0,
		TLSNextProto:   nil,
		ConnState:      nil,
		ErrorLog:       nil,
		BaseContext:    nil,
		ConnContext:    nil,
	}
	s.SetKeepAlivesEnabled(false)
	if err := s.Serve(l); err != nil {
		log.Panicf("Error serving: %v", err.Error())
	}

}

func handle(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second * 1)
	w.Write([]byte("hello"))
	log.Println("finish")

}
