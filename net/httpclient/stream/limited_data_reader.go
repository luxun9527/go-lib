package stream

import "io"

var (
	_data = make([]byte, 1024*800)
)

// LimitedDataReader generates a specific amount of data but restricts the read size.
type LimitedDataReader struct {
	remaining   int64
	perReadSize int64
}

func (d *LimitedDataReader) Read(p []byte) (n int, err error) {
	if d.remaining <= 0 {
		return 0, io.EOF
	}

	// Determine the size of the chunk to generate.
	chunkSize := int64(len(p))
	if chunkSize > d.perReadSize {
		chunkSize = d.perReadSize
	}
	if chunkSize > d.remaining {
		chunkSize = d.remaining
	}

	// Fill the buffer with generated data (for example, 'A' characters).
	copy(p[:chunkSize], _data)
	d.remaining -= chunkSize
	return int(chunkSize), nil
}

func NewLimitedDataReader(size int64, perReadSize int64) *LimitedDataReader {
	if perReadSize > int64(len(_data)) {
		perReadSize = int64(len(_data))
	}
	return &LimitedDataReader{remaining: size, perReadSize: perReadSize}
}
