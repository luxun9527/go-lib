package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	token      = "token=cimf8epr01qlsedscmvgcimf8epr01qlsedscn00"
	HttpTarget = "www.finnhub.io"
	WsTarget   = "ws.finnhub.io"
)

func main() {
	if err := http.ListenAndServe(":52023", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			//为ws请求
			if req.Header.Get("Upgrade") != "" {
				req.URL.Host = WsTarget
				req.Host = WsTarget
			} else {
				req.URL.Host = HttpTarget
				req.Host = HttpTarget
			}
			req.URL.Scheme = "https"
			delimiter := "&"
			if req.URL.RawQuery == "" {
				delimiter = ""
			}
			req.URL.RawQuery += delimiter + token
			req.RequestURI += delimiter + token
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})); err != nil {
		log.Fatal(err)
	}
}
