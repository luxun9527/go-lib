package main

import (
	"github.com/go-lark/lark"
	"log"
)

func main() {
	Json()
}

func Json() {
	b := lark.NewNotificationBot("https://open.feishu.cn/open-apis/bot/v2/hook/71f86ea6-ab9a-4645-b40b-1be00ffe910a")
	builder := lark.NewCardBuilder()
	card := builder.Card(builder.Markdown("```json \n{\"code\":200}")).Yellow().Title("test")
	msg := lark.NewMsgBuffer(lark.MsgInteractive)
	om := msg.Card(card.String()).Build()
	resp, err := b.PostNotificationV2(om)
	if err != nil {
		log.Printf("send message failed: %v", err)
	}
	log.Printf("send message response: %v", resp)
}
func String() {
	// 替换自己的机器人地址
	b := lark.NewNotificationBot("https://open.feishu.cn/open-apis/bot/v2/hook/71f86ea6-ab9a-4645-b40b-1be00ffe910a")
	builder := lark.NewCardBuilder()
	card := builder.Card(builder.Markdown("```json \n{\"code\":200}")).Yellow().Title("test")
	msg := lark.NewMsgBuffer(lark.MsgInteractive)
	om := msg.Card(card.String()).Build()
	resp, err := b.PostNotificationV2(om)
	if err != nil {
		log.Printf("send message failed: %v", err)
	}
	log.Printf("send message response: %v", resp)

}
