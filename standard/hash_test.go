package standard

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"testing"
)

func TestHash(t *testing.T) {
	c := sha1.New()
	c.Write([]byte("input"))
	bytes := c.Sum(nil)
	log.Println( hex.EncodeToString(bytes))

}
