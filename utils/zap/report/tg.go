package report

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"time"
)

type TgWriter struct {
	bot    *tgbotapi.BotAPI
	chatId int64
}

func NewTgWriter(token string, chatId int64) *TgWriter {
	transport := &http.Transport{
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     50,
		IdleConnTimeout:     30 * time.Second,
	}
	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, &http.Client{
		Transport: transport,
	})
	if err != nil {
		log.Panicf("init tg bot failed, err:%v", err)
	}
	return &TgWriter{bot: bot, chatId: chatId}
}
func (l *TgWriter) Write(p []byte) (n int, err error) {
	msg := tgbotapi.NewMessage(l.chatId, string(p))
	if _, err = l.bot.Send(msg); err != nil {
		return 0, err
	}
	return len(p), nil
}
func (l *TgWriter) Sync() error {
	return nil
}
