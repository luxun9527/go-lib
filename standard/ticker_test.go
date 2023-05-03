package standard

import (
	"log"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second * 2)
	for {

		t1 := <-ticker.C
		//time.Sleep(time.Second * 3)
		log.Println(t1)

		log.Println("test")
		ticker.Reset(time.Second)
	}
}
