package goutil

import (
	"github.com/gookit/goutil/netutil"
	"github.com/zeromicro/go-zero/core/netx"
	"log"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	{
		host := netutil.InternalIPv4()
		log.Printf("host:%s", host)
	}
	{
		host := netutil.InternalIPv1()
		log.Printf("host:%s", host)
	}
	{
		host := netx.InternalIp()
		log.Printf("host:%s", host)
	}
}
