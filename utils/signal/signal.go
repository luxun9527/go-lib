package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select  {
		case sig:= <-c: {
			log.Printf("Got %s signal. Aborting...", sig)
		}
		}
	}()

	select {

	}
}