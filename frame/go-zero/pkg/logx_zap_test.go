package pkg

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestZapLog(t *testing.T) {
	w := NewZapWriter(&ZapLogger{
		Level:         "info",
		Stacktrace:    false,
		AddCaller:     true,
		Mode:          "console",
		FileName:      "",
		ErrorFileName: "",
		MaxSize:       0,
		MaxAge:        0,
		MaxBackup:     0,
		Async:         false,
		Json:          false,
		Compress:      false,
		options:       nil,
	})
	logx.SetWriter(w)
	logx.Info("test1 Info")
	logx.Debug("test debug")
	logx.Error("test Error")
}
