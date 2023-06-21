package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	r := gin.New()
	r.GET("/echo", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		conn.SetPingHandler(func(c string) error {
			log.Println("data = ", c)
			return nil
		})
		conn.SetCloseHandler(func(code int, text string) error {
			log.Printf("code = %v text = %v", code, text)
			return nil
		})
		//	conn.SetReadDeadline(time.Now().Add(time.Second * 20))
		go func() {
			for {
				_, data, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						log.Println("websocket.CloseNormalClosure message")
					}
					log.Println("read message failed", err)
					break
				}
				log.Println("data", string(data))
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Println("write message failed", err)
					break
				}
			}
		}()
	})
	r.Run(":8989")
}
