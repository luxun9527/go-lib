package report

import (
	"github.com/go-lark/lark"
)

type LarkWriter struct {
	bot *lark.Bot
}

func NewLarkWriter(token string) *LarkWriter {
	bot := lark.NewNotificationBot(token)
	return &LarkWriter{
		bot: bot,
	}
}

func (l *LarkWriter) Write(p []byte) (n int, err error) {
	_, err = l.bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text(string(p)).Build())
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
func (l *LarkWriter) Sync() error {
	return nil
}
