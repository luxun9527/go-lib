package logger

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestZap(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    1,
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
	}
	l := config.Build()
	l.Debug("this a debug level log", zap.Any("test", "t"))
	l.Info("this a info level log", zap.Any("test", "t"))
	l.Warn("this a warn level log", zap.Any("test", "t"))
	l.Error("this a error level log", zap.Any("test", "t"))
//	l.Sync()
    time.Sleep(time.Second*2)
	l.Panic("this a panic level log", zap.Any("test", "t"))

}
func TestZapAsync(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    1,
		Mode:          "console",
		FileName:      "",
		ErrorFileName: "",
		MaxSize:       0,
		MaxAge:        0,
		MaxBackup:     0,
		Async:         true,
		Json:          true,
		Compress:      false,
		options:       nil,
	}
	l := config.Build()
	l.Debug("this a debug level log", zap.Any("test", "t"))
	l.Info("this a info level log", zap.Any("test", "t"))
	l.Warn("this a warn level log", zap.Any("test", "t"))
	l.Error("this a error level log", zap.Any("test", "t"))
	//	l.Sync()
	//l.Panic("this a panic level log", zap.Any("test", "t"))
	time.Sleep(time.Second*2)
	l.Sync()

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
		MaxBackup:     10,
		Async:         true,
		Json:          false,
		Compress:      false,
		options:       nil,
	}
	l := config.Build()
	for i := 0; i < 13000000; i++ {
		l.Debug("this a debug level log", zap.Any("test", "t"))
		l.Info("this a info level log", zap.Any("test", "t"))
		l.Warn("this a warn level log", zap.Any("test", "t"))
		l.Error("this a error level log", zap.Any("test", "t"))

	}
	l.Sync()

}
