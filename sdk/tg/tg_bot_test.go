package tg

import (
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestInit(t *testing.T) {
	proxyURL, _ := url.Parse("http://127.0.0.1:7890") // 替换为你的Clash代理地址和端口

	transport := &http.Transport{
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     50,
		IdleConnTimeout:     30 * time.Second,
		Proxy:               http.ProxyURL(proxyURL),
	}
	bot, err := tgbotapi.NewBotAPIWithClient("6873378295:AAEisacb2Q1Re2FPiyoS8f2hkn4uckMQ3ck", tgbotapi.APIEndpoint, &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	})
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
