package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"testing"
)

func TestCookie(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/set", func(writer http.ResponseWriter, request *http.Request) {
		cookie1 := &http.Cookie{
			Name:       "name",
			Value:      "zhangsan",
			Path:       "",
			Domain:     "",
			RawExpires: "",
			MaxAge:     3600,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		}
		http.SetCookie(writer, cookie1)
		cookie2 := &http.Cookie{
			Name:       "age",
			Value:      "12",
			Path:       "",
			Domain:     "",
			RawExpires: "",
			MaxAge:     3600,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		}
		http.SetCookie(writer, cookie2)
	})
	r.HandleFunc("/read", func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("name")
		if err != nil {
			log.Println(err)
		}
		log.Println(cookie.Value)
	})
	log.Fatal(http.ListenAndServe(":8082", r))
}
