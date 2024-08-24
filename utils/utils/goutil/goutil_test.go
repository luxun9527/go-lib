package goutil

import (
	"fmt"
	"net"
	"testing"
)

func TestGoUtils(t *testing.T) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, addr := range addrs {
		// 检查地址类型并跳过回环地址
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// 仅打印IPv4地址
			if ipNet.IP.To4() != nil {
				fmt.Println("Pod IP Address:", ipNet.IP.String())
			}
		}
	}
}
