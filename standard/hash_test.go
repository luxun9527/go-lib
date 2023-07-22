package standard

import (
	"fmt"
	"hash/fnv"
	"testing"
)

func TestHash(t *testing.T) {
	f := fnv.New32()
	f.Write([]byte("hello1"))
	fmt.Println(f.Sum32())
}
