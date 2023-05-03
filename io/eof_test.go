package io

import (
	"strings"
	"testing"
)

func TestEOF(t *testing.T) {
	reader := strings.NewReader("ffffffssadffsa")
	buf := make([]byte, 10)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			t.Log(err)
			return
		}
		t.Log(string(buf[:n]))
	}
}
