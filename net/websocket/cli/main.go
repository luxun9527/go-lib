package main

import (
	"flag"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
)

var (
	currentReceived atomic.Int64
	connected       atomic.Int64
	lastTime        int64
	lastReceived    int64
)

func main() {
	var url string
	var count int64

	flag.StringVar(&url, "url", "ws://47.113.223.16:9993/ws", "url")
	flag.Int64Var(&count, "count", 5000, "count")
	flag.Parse()
	for i := int64(0); i < count; i++ {
		time.Sleep(time.Microsecond)
		go connect(url)
	}
	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now().UnixMilli()
			cr := currentReceived.Load()
			s := now - lastTime
			receivedLastSecond := (cr - lastReceived) / s
			lastReceived = cr
			lastTime = now
			log.Printf("当前连接数%v当前收到 %v条 平均每毫秒收到 %v条\n", connected.Load(), cr, receivedLastSecond)
		}

	}()
	time.Sleep(time.Hour)
}
func connect(url string) {
	dialer := websocket.Dialer{
		EnableCompression: true, // 启用压缩
		Subprotocols:      []string{"permessage-deflate"},
	}

	//c := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Printf("connect failed err = %v", err.Error())
		return
	}
	ReadMessage(conn)
}
func ReadMessage(conn *websocket.Conn) {
	conn.EnableWriteCompression(true)
	conn.SetCompressionLevel(9)
	conn.WriteJSON(map[string]interface{}{"code": 1, "topic": "test"})
	connected.Inc()
	for {
		dataType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("data %v err %v %v", string(data), err, dataType)
			connected.Dec()
			return
		}
		currentReceived.Inc()

	}
}
