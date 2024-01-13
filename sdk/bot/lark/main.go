package main

import (
	"github.com/go-lark/lark"
	"log"
)

func main() {
	bot := lark.NewNotificationBot("xxxxx")
	result, err := bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text("hello, wolrd").Build())
	if err != nil {
		log.Printf("post message failed, err:%v", err)
	}
	log.Println(result)
}
