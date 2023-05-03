package io

import (
	"bufio"
	"log"
	"testing"
)

type Wr struct{}

func (*Wr) Write(p []byte) (n int, err error) {
	log.Println(string(p))
	return 0, err
}

func TestBufio(t *testing.T) {
	var r Wr
	wr := bufio.NewWriterSize(&r, 2024)
	for i := 0; i < 1000; i++ {
		wr.Write([]byte{1, 2, 3, 4})
	}
}
