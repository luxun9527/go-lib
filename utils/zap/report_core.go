package logger

import (
	"github.com/go-lark/lark"
	"sync"
	"time"
)

type LarkWriter struct {
	bot *lark.Bot
}

func (l *LarkWriter) Write(p []byte) (n int, err error) {
	_, err = l.bot.PostNotificationV2(lark.NewMsgBuffer(lark.MsgText).Text(string(p)).Build())
	if err != nil {
		return 0, err
	}
	return len(p), err
}
func (l *LarkWriter) Sync() error {
	return nil
}

func newLarkWrite() *LarkWriter {
	bot := lark.NewNotificationBot("xxxxx")
	return &LarkWriter{bot: bot}
}

func NewLarkWriterBuffer() *LarkWriterBuffer {
	return &LarkWriterBuffer{
		buf: make([]string, 0, 15),
		w:   newLarkWrite(),
	}
}

type LarkWriterBuffer struct {
	lock sync.Mutex
	buf  []string
	w    *LarkWriter
}

func (l *LarkWriterBuffer) Start() {
	for {
		time.Sleep(time.Second)
		l.Sync()
	}
}

func (l *LarkWriterBuffer) Write(p []byte) (n int, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	// 缓存太多写入了，则抛弃原来的
	if len(l.buf) >= 10 {
		//
		return
	}
	l.buf = append(l.buf, string(p))
	return len(p), err
}

func (l *LarkWriterBuffer) Sync() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, v := range l.buf {
		l.w.Write([]byte(v))
	}
	l.buf = make([]string, 0, 10)
	return nil
}
