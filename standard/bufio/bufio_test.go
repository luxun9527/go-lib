package tools

import (
	"io"
	"testing"
)

//io.ReadFull不足的情况下回退N步
func TestBufioRead(t *testing.T) {
	var customReader CustomReader
	reader := NewReaderSize(customReader, 1024)
	buf := make([]byte, 2)
	n, err := io.ReadFull(reader, buf)
	if err != io.ErrShortWrite {
		reader.GoBackN(n)
	}
}

type CustomReader struct {
}

var f = true

func (CustomReader) Read(p []byte) (n int, err error) {
	if f {
		p[0] = 1
		f = false
		return 1, err
	}
	return 0, io.EOF
}
