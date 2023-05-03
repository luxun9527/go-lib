package io

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestTruncate(t *testing.T) {
	fs, err := os.OpenFile("test.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("err", err)
		return
	}
	go func() {
		for {
			time.Sleep(time.Second)
			fs.WriteString("aaaaa\n")
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second)
			if err := fs.Truncate(0); err != nil {
				log.Println("err", err)
			}
		}
	}()
	time.Sleep(time.Hour)
}
