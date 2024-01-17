package standard

import (
	"log"
	"testing"
)

func TestPointer(t *testing.T) {
	var i = 1
	log.Printf("%p\n", &i)
	p := &i
	log.Printf("%p\n", p)

}
