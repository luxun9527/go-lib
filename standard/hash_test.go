package standard

import (
	"go.uber.org/zap"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	v1 := zap.AtomicLevel{}
	v2 := zap.NewAtomicLevelAt(1)
	log.Printf("%v", v1 == v2)

}
