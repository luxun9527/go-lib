package main

import (
	"github.com/go-lark/lark"
	"log"
)

func main() {
	bot :=lark.NewNotificationBot("https://open.larksuite.com/open-apis/bot/v2/hook/4d969fd8-c2ab-494a-a698-49acb7fcbb98")
	result, err := bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text("hello, wolrd").Build())
	if err!=nil{
		log.Printf("post message failed, err:%v", err)
	}
	log.Println(result)
}


