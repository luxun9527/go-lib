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

func TestZapConsole(t *testing.T) {
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
	l.Debug("this a debug level log", zap.Any("test", "debug"))
	l.Info("this a info level log", zap.Any("test", "info"))
	l.Warn("this a warn level log", zap.Any("test", "warn"))
	l.Error("this a error level log", zap.Any("test", "error"))
	l.Panic("this a panic level log", zap.Any("test", "panic"))
}
func TestZapFile(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    0,
		Mode:          "file",
		FileName:      "stdout.log",
		ErrorFileName: "stderr.log",
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
	/*
		2024-01-08-22:33:31	DEBUG	E:/demoproject/go-lib/utils/zap/zap_test.go:33	this a debug level log	{"test": "debug"}

		2024-01-08-22:33:31	INFO	E:/demoproject/go-lib/utils/zap/zap_test.go:34	this a info level log	{"test": "info"}
		2024-01-08-22:33:31	WARN	E:/demoproject/go-lib/utils/zap/zap_test.go:35	this a warn level log	{"test": "warn"}
		2024-01-08-22:33:31	ERROR	E:/demoproject/go-lib/utils/zap/zap_test.go:36	this a error level log	{"test": "error"}

		2024-01-08-22:33:31	PANIC	E:/demoproject/go-lib/utils/zap/zap_test.go:37	this a panic level log	{"test": "panic"}
	*/
}

func TestZapAsync(t *testing.T) {
	config := Config{
		Level:         "debug",
		Stacktrace:    true,
		AddCaller:     true,
		CallerShip:    0,
		Mode:          "console",
		FileName:      "std.log",
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

	config.UpdateLevel(zapcore.ErrorLevel)
	l.Debug("this a debug level log", zap.Any("test", "t"))
	l.Info("this a info level log", zap.Any("test", "t"))
	l.Warn("this a warn level log", zap.Any("test", "t"))
	l.Error("this a error level log", zap.Any("test", "t"))
}
