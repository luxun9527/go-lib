package main

import (
	"github.com/gorilla/websocket"
	"log"
	"testing"
)

func TestClient(t *testing.T) {
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token=cimf8epr01qlsedscmvgcimf8epr01qlsedscn00", nil)
	if err != nil {
		log.Panic(err)
	}
	if err := c.WriteMessage(websocket.TextMessage, []byte(`{"type":"subscribe","symbol":"BINANCE:BTCUSDT"}`)); err != nil {
		log.Panic(err)
	}
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			log.Panic(err)
		}
		log.Println("Received:", string(data))
	}

}
