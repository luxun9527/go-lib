package http

import (
	"log"
	"net/http"
	"testing"
)

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Print("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareTwo")
		if r.URL.Path == "/foo" {
			return
		}

		next.ServeHTTP(w, r)
		log.Print("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func TestHttpGoroutine(t *testing.T) {
	go func() {
		http.ListenAndServe("0.0.0.0:8899", nil)
	}()
	if err := http.ListenAndServe(":8888", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		select {
		case <-request.Context().Done():
			log.Println("request canceled")
		}
	})); err != nil {
		log.Fatal(err)
	}
}

func TestHttp(t *testing.T) {
	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(final)
	mux.Handle("/", middlewareOne(middlewareTwo(finalHandler)))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
