package main

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func main() {
	addr := "ws://192.168.2.159:9995/ws"

	for i := 0; i < 1000; i++ {
		Connect(addr)
	}
	time.Sleep(time.Second * 30)
	for i := 0; i < 1000; i++ {
		Connect(addr)
	}
	time.Sleep(time.Hour)
}
func Connect(addr string) {
	time.Sleep(time.Microsecond * 2)
	go func() {
		c, _, err := websocket.DefaultDialer.Dial(addr, nil)
		if err != nil {
			log.Println("err", err)
			return
		}
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

}
