package main

import (
	"github.com/gorilla/websocket"
	"github.com/luxun9527/zlog"
	"testing"
	"time"
)

func TestRb(t *testing.T) {

	//data, _ := os.ReadFile("D:\\weardu\\app\\ai\\aigc\\etc\\prompt.yaml")
	//var m map[string]string
	//if err := yaml.Unmarshal(data, &m); err != nil {
	//	zlog.Panicf("unmarshal failed: %v", err)
	//}
	//for k, v := range m {
	//	zlog.Infof("key=%v,value=%v", k, v)
	//}
	conn, _, err := websocket.DefaultDialer.Dial("wss://test.api.weardu.com/aiws/v1/ai/ws?source=1&token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHRyYSI6eyJ1aWQiOjcsInVzZXJuYW1lIjoiIiwicGxhdGZvcm0iOiIifSwiZXhwIjoxNzQ5NTIzOTExfQ.qN4Przwdq49eCAapTQzGWwrVTr0GgSxPl5mmkptaobk&mode=chat&srcLang=auto&targetLang=en", nil)
	if err != nil {
		zlog.Panicf("err=%v", err)
	}
	defer conn.Close()
	conn.SetPongHandler(func(appData string) error {
		zlog.Infof("appData=%v", appData)
		return nil
	})
	go func() {
		for {
			time.Sleep(time.Second * 10)
			if err := conn.WriteMessage(websocket.PingMessage, []byte("test")); err != nil {
				zlog.Panicf("err=%v", err)
			}
		}
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			zlog.Errorf("err=%v", err)
			return
		}

	}
}
