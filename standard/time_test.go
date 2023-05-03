package standard

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestJsonTime(t *testing.T) {
	type TimeFormat struct {
		UpdateTime time.Time `json:"updateTime"`
	}
	//data := `{"updateTime":"2023-03-31T13:03:57+08:00"}`
	//data := `{"updateTime":"2023-03-28T02:48:17.000000000Z"}`
	data := `{"updateTime":"2023-03-28T15:52:53Z"}`
	var tf TimeFormat
	if err := json.Unmarshal([]byte(data), &tf); err != nil {
		log.Fatal(err)
	}
	log.Println(tf)

}
