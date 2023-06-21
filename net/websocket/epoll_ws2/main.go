package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	gws "github.com/gobwas/ws"
	"log"
	"net"
	"time"
)

func main() {
	engine := gin.New()
	Start()

	engine.GET("/ws", Connect)
	engine.Run(":8989")
}

type Connection struct {
	buf  *bytes.Buffer
	conn net.Conn
	id   int
}

func NewConnection(conn net.Conn) *Connection {
	fd := WebsocketFD(conn)
	connection := &Connection{
		buf:  bytes.NewBuffer(make([]byte, 0, 100)),
		conn: conn,
		id:   fd,
	}
	return connection
}

func Connect(c *gin.Context) {
	var httpUpgrade gws.HTTPUpgrader
	conn, _, _, err := httpUpgrade.Upgrade(c.Request, c.Writer)
	if err != nil {
		return
	}
	if err := Epoller.Add(conn); err != nil {
		log.Println("err", err)
		return
	}
	time.Sleep(time.Second * 3)
	conn.Close()
}
