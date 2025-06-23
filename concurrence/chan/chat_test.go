package _chan

import (
	"testing"
	"time"
)

func TestChan(t *testing.T) {

	ch := make(chan int)

	go func() {
		<-ch
		t.Log("ok")
		<-ch
		t.Log("ok")

	}()

	ch <- 1
	time.Sleep(time.Second * 2)
}
