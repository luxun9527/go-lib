package goutil

import (
	"github.com/gookit/goutil/netutil"
	"log"
	"testing"
)

func TestGoUtils(t *testing.T) {

	log.Println(netutil.GetLocalIPs())
	print(netutil.IPv4())

}
