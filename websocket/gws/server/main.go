package main

import (
	"github.com/gin-gonic/gin"
	gws "github.com/gobwas/ws"
	"log"
)

func main() {
	engine := gin.New()
	engine.GET("/ws", Connect)
	engine.Run(":8989")
}
func Connect(c *gin.Context) {
	var httpUpgrade gws.HTTPUpgrader
	conn, _, _, err := httpUpgrade.Upgrade(c.Request, c.Writer)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	//设置超时时间
	//for {
	frame, err := gws.ReadFrame(conn)
	if err != nil {
		log.Println("err", err)
		return
	}
	log.Println(string(frame.Payload))

	//}
}
