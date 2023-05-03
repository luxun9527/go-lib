package standard

import (
	"encoding/json"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Data []byte `json:"data"`
	}

	u := User{
		Name: "aa",
		Data: []byte("abc"),
	}
	r, _ := json.Marshal(u)
	log.Println(string(r))
}
