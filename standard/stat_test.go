package standard

import (
	"log"
	"os"
	"testing"
)

//go test -v stat_test.go -test.run TestStat
func TestStat(t *testing.T) {
	_, err := os.Stat("")
	if err != nil {
		log.Println("err", err)
	}

}
