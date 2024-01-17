package report

import (
	"fmt"
	wxworkbot "github.com/vimsucks/wxwork-bot-go"
)

type WxWriter struct {
	bot *wxworkbot.WxWorkBot
}

func NewWxWriter(token string) *WxWriter {
	bot := wxworkbot.New(token)
	return &WxWriter{
		bot: bot,
	}
}

/*
```json

```
*/

func (l *WxWriter) Write(p []byte) (n int, err error) {
	markdown := wxworkbot.Markdown{
		Content: fmt.Sprintf("```json\n%s\"\n````", string(p)),
	}
	if err = l.bot.Send(markdown); err != nil {
		return 0, err
	}
	return len(p), nil
}
func (l *WxWriter) Sync() error {
	return nil
}
