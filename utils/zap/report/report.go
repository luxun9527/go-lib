package report

import (
	"bufio"
	"errors"
	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
	"time"
)

/*
日志上报，告知异常，具体错误还是通过日志排查
*/

var (
	ReportChanFull = errors.New("report chan full")
)

type ImType string

const (
	Tg   ImType = "tg"
	Wx          = "wx"
	Lark        = "lark"
)

type ImConfig struct {
	//上报的类型，目前支持：wx lark tg
	Type string `json:"type"`
	//lark填webhook tg wx填token
	Token string `json:"token"`
	//tg的chatid
	ChatID string `json:",optional"`
	//日志刷新的频率 单位秒
	FlushSec string `json:",default=3"`
}

func NewWriteSyncer(c ImConfig) zapcore.WriteSyncer {
	switch ImType(c.Type) {
	case Wx:
	case Lark:
	case Tg:
	default:
		log.Panicf("unsupported report type:%s", c.Type)
	}
	return nil
}

func NewImWriterBuffer(w zapcore.WriteSyncer) *ImWriterBuffer {
	return &ImWriterBuffer{
		buf: bufio.NewWriterSize(w, 1024*10),
		w:   w,
	}
}

type ImWriterBuffer struct {
	lock     sync.Mutex
	buf      *bufio.Writer
	w        zapcore.WriteSyncer
	sendChan chan []byte
	count    atomic.Int32
	flushSec int64
}

const _maxCount = 10

func (l *ImWriterBuffer) Start() {
	t := time.NewTicker(time.Second * time.Duration(l.flushSec))
	for {
		select {
		case <-t.C:
		case data := <-l.sendChan:
			//希望缓存永远不会达到阈值，缓存达到阈值触发刷新，可能发送的数据不全
			if l.count.Load() > _maxCount {
				l.buf.Flush()
			}

			l.buf.Write(data)
		}

	}
}

func (l *ImWriterBuffer) Write(p []byte) (n int, err error) {

	select {
	case l.sendChan <- p:
		l.count.Add(1)
	default:
		return 0, ReportChanFull
	}

	return len(p), err
}

func (l *ImWriterBuffer) Sync() error {
	return l.buf.Flush()
}
