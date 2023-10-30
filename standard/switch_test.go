package standard

import (
	"log"
	"testing"
)

func TestSwitch(t *testing.T) {
	var i,j int
	switch  {
	case i==0:
		log.Println(1)
	case j==0:
		log.Println(2)
	}

}
