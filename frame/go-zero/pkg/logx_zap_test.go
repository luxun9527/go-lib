package pkg

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestZapLog(t *testing.T) {
	logx.Error("test1")
	w := NewZapWriter(&Config{
		Level:         "info",
		Stacktrace:    false,
		AddCaller:     true,
		Mode:          "console",
		FileName:      "test1.log",
		ErrorFileName: "",
		MaxSize:       0,
		MaxAge:        0,
		MaxBackup:     0,
		Async:         false,
		Json:          false,
		Compress:      false,
		options:       nil,
		CallerShip:    3,
	})
	logx.SetWriter(w)
	logx.Infow("test1 Info")
	logx.Debug("test debug")
	logx.Error("test Error")
	logx.Slow("test")
	logx.Severe("")
}
