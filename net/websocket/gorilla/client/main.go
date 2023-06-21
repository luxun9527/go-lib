package main

import (
	"github.com/gorilla/websocket"

	"log"
	"time"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:9090/echo", nil)
	if err != nil {
		log.Println("err", err)
		return
	}
	go func() {
		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				log.Println("err", err)
				return
			}
			log.Println("data", string(data))
		}

	}()
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Second)
		if err := c.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
			log.Println("err", err)
			return
		}
		if i == 10 {
			c.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(time.Second*3))
		}
		if i == 20 {
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*3))
		}
	}
}
