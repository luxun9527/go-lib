package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)
//https://www.cnblogs.com/zy666/p/10504272.html
func main() {
	//go env -w GOOS=linux
	c := make(chan os.Signal)
	//监听ctrl+c和kill命令
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-c:
		log.Printf("Got %s signal. Aborting...", sig)

	}

}
