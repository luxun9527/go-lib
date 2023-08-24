package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	// ListenAndServeTLS always returns a non-nil error.
	err := http.ListenAndServeTLS(":10443", "E:\\demoproject\\go-lib\\net\\https\\cert.pem", "E:\\demoproject\\go-lib\\net\\https\\key.pem", nil)
	log.Fatal(err)
}
