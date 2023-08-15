package main

import (
	"fmt"
	"github.com/boj/redistore"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var (
	store      = sessions.NewFilesystemStore("./", securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
	redisStore *redistore.RediStore
)

func set(w http.ResponseWriter, r *http.Request) {
	session, _ := redisStore.Get(r, "user")
	session.Values["name"] = "dj"
	session.Values["age"] = 18
	err := redisStore.Save(r, w, session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Hello World")
}

func read(w http.ResponseWriter, r *http.Request) {
	session, _ := redisStore.Get(r, "user")

	fmt.Fprintf(w, "name:%s age:%d\n", session.Values["name"], session.Values["age"])
}

func main() {
	InitRedisStore()
	r := mux.NewRouter()
	r.HandleFunc("/set", set)
	r.HandleFunc("/read", read)
	log.Fatal(http.ListenAndServe(":8081", r))
}
func InitRedisStore() {
	var err error
	redisStore, err = redistore.NewRediStore(10, "tcp", "192.168.2.99:6379", "", []byte("secret-key"))
	if err != nil {
		log.Fatal(err)
	}
}
