package main

import (
	"github.com/vimsucks/wxwork-bot-go"
	"log"
)

func main() {
	const webhookAddr = "xxx"
	bot := wxworkbot.New(webhookAddr)
	// or Markdown, Image, News

	// 文本消息
	text := wxworkbot.Text{
		Content:             "Hello World",
		MentionedList:       []string{"foo", "bar"},
		MentionedMobileList: []string{"@all"},
	}
	err := bot.Send(text)
	if err != nil {
		log.Fatal(err)
	}

	// Markdown 消息
	markdown := wxworkbot.Markdown{
		Content: "# 测试",
	}
	err = bot.Send(markdown)
	if err != nil {
		log.Fatal(err)
	}

}
