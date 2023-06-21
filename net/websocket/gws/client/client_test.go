package client

import (
	"bytes"
	"context"
	"github.com/gobwas/ws/wsflate"
	"log"
	"testing"
)
import (
	gws "github.com/gobwas/ws"
)

func TestClient(t *testing.T) {
	dialer := gws.DefaultDialer
	conn, _, _, err := dialer.Dial(context.Background(), "ws://192.168.2.99:9992/ws")
	if err != nil {
		log.Fatalln("err", err)
	}

	message := NewMessage(gws.OpText, []byte("abcdefjhigklmnopqrusasdflsdfasdl"))
	b, _ := message.ToWebSocketFrame()
	for i := 0; i < 10; i++ {
		conn.Write(b)
	}
	select {}
}

type Message struct {
	messageType gws.OpCode
	data        []byte
}

func NewMessage(code gws.OpCode, data []byte) Message {
	return Message{
		messageType: code,
		data:        data,
	}
}

func (m Message) ToWebSocketFrame() ([]byte, error) {
	var res bytes.Buffer
	frame := gws.NewFrame(m.messageType, true, m.data)
	if err := gws.WriteFrame(&res, frame); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
func (m Message) ToCompressFrame() ([]byte, error) {
	var res bytes.Buffer
	frame := gws.NewFrame(m.messageType, true, m.data)
	compressFrame, err := wsflate.CompressFrame(frame)
	if err != nil {
		return nil, err
	}
	//res.Write()
	if err := gws.WriteFrame(&res, compressFrame); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}
