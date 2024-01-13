package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestZapReport(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    0,
		Mode:          "console",
		FileName:      "",
		ErrorFileName: "",
		MaxSize:       0,
		MaxAge:        0,
		MaxBackup:     0,
		Async:         false,
		Json:          true,
		Compress:      false,
		IsReport: true,
		options:       nil,
	}
	l := config.Build()

	l.Debug("this a debug level log", zap.Any("test", "debug"))
	l.Info("this a info level log", zap.Any("test", "info"))
	l.Warn("this a warn level log", zap.Any("test", "warn"))
	l.Error("this a error level log", zap.Any("test", "error"))
	l.Panic("this a panic level log", zap.Any("test", "panic"))
}