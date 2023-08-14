package logger

import (
	"go.uber.org/zap"
	"testing"
)

func TestZap(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    0,
		Mode:          "console",
		FileName:      "0",
		ErrorFileName: "",
		MaxSize:       0,
		MaxAge:        0,
		MaxBackup:     0,
		Async:         true,
		Json:          false,
		Compress:      false,
		options:       nil,
	}
	l := config.Build()
	l.Debug("this a debug level log", zap.Any("test", "t"))
	l.Info("this a info level log", zap.Any("test", "t"))
	l.Warn("this a warn level log", zap.Any("test", "t"))
	l.Error("this a error level log", zap.Any("test", "t"))
	l.Sync()
	l.Panic("this a panic level log", zap.Any("test", "t"))

}
func TestFileLog(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    0,
		Mode:          "file",
		FileName:      "std.log",
		ErrorFileName: "err.log",
		MaxSize:       1,
		MaxAge:        1,
		MaxBackup:     1,
		Async:         true,
		Json:          false,
		Compress:      false,
		options:       nil,
	}
	l := config.Build()
	for i := 0; i < 1000000; i++ {
		l.Debug("this a debug level log", zap.Any("test", "t"))
		l.Info("this a info level log", zap.Any("test", "t"))
		l.Warn("this a warn level log", zap.Any("test", "t"))
		l.Error("this a error level log", zap.Any("test", "t"))

	}
	l.Sync()

}
