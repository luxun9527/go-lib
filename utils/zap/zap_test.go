package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

/*
test 主要是作为示例

*/

func TestZap(t *testing.T) {
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
		Json:          false,
		Compress:      false,
		options:       nil,
	}
	l := config.Build()
	l.Debug("this a debug level log", zap.Any("test", "t"))
	l.Info("this a info level log", zap.Any("test", "t"))
	l.Warn("this a warn level log", zap.Any("test", "t"))
	l.Error("this a error level log", zap.Any("test", "t"))
	l.Panic("this a panic level log", zap.Any("test", "t"))
}
func TestZapAsync(t *testing.T) {
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
		Async:         true,
		Json:          false,
		Compress:      false,
		options:       nil,
	}
	l := config.Build()
	l.Debug("this a debug level log")
	l.Info("this a info level log")
	l.Warn("this a warn level log")
	time.Sleep(time.Second * 2)
	l.Sync()
}
func TestFileLog(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    0,
		Mode:          "file",
		FileName:      "stdout.log",
		ErrorFileName: "stderr.log",
		MaxSize:       1,
		MaxAge:        1,
		MaxBackup:     10,
		Async:         true,
		Json:          true,
		Compress:      true,
		options:       nil,
	}
	l := config.Build()
	for i := 0; i < 150000; i++ {
		l.Debug("this a debug level log", zap.Any("test", "t"))
		l.Info("this a info level log", zap.Any("test", "t"))
		l.Warn("this a warn level log", zap.Any("test", "t"))
		l.Error("this a error level log", zap.Any("test", "t"))
	}

	l.Sync()

}
func TestChangeOnRuntime(t *testing.T) {
	//在运行时改变日志等级
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
		Json:          false,
		Compress:      false,
		options:       nil,
	}
	l := config.Build()
	l.Debug("this a debug level log", zap.Any("test", "t"))
	l.Info("this a info level log", zap.Any("test", "t"))
	l.Warn("this a warn level log", zap.Any("test", "t"))
	l.Error("this a error level log", zap.Any("test", "t"))

	config.ChangeLevel(zapcore.ErrorLevel)
	l.Debug("this a debug level log", zap.Any("test", "t"))
	l.Info("this a info level log", zap.Any("test", "t"))
	l.Warn("this a warn level log", zap.Any("test", "t"))
	l.Error("this a error level log", zap.Any("test", "t"))
}
