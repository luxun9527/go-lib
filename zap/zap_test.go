package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"testing"

	"time"
)

func TestBuildZap(t *testing.T) {
	loggConfig := LoggerConfig{
		Level:              "debug",
		Stacktrace:         true,
		AddCaller:          true,
		FileName:           "test.log",
		EnableWarnRedirect: true,
		WarnFileName:       "warn.log",
		MaxSize:            1,
		MaxAge:             1,
		MaxBackup:          5,
		Async:              false,
		Json:               false,
		Mode:               "file",
	}
	logger := loggConfig.Build()
	logger.Debug("test Debug level")
	logger.Info("test Info")
	logger.Warn("test Warn")
	logger.Error("test Error")
	logger.Panic("test Panic")
}

func (lc *LoggerConfig) ParseLevel() zapcore.Level {
	level, err := zapcore.ParseLevel(lc.Level)
	if err != nil {
		log.Panicf("init level failed level %s err %v", lc.Level, err)
	}
	return level
}

type LoggerConfig struct {
	//日志级别 debug info warn panic
	Level string
	//panic，是否显示堆栈 panic级别的日志输出堆栈信息。
	Stacktrace bool
	//添加调用者信息
	AddCaller bool
	//输出到哪里标准输出console,还是文件file
	Mode string
	//文件名称加路径
	FileName string
	//是否开启warn级别以上的日志重定向，在console,warn级别的日志输出到标准错误输出中，在file模式中warn输出到WarnFileName中，其他的配置相同。
	EnableWarnRedirect bool
	//warn 级别的日志输出到不同的地方
	WarnFileName string
	// 日志轮转大小，单位MB，默认500MB
	MaxSize int
	//日志轮转最大时间，单位day，默认1 day
	MaxAge int
	//日志最大保留的个数
	MaxBackup int
	//异步日志 待实现
	Async bool
	//是否 输出json格式的数据，JSON格式相对于console格式，不方便阅读，但是对机器更加友好
	//最佳实践，在开发的时候json为false,mode为console,测试部署阶段，EnableWarnRedirect为true,mode为file,json为TRUE
	Json bool
	//是否日志压缩
	Compress bool
	options  []zap.Option
}

func (lc *LoggerConfig) Build() *zap.Logger {

	var (
		ws      zapcore.WriteSyncer
		warnWs  zapcore.WriteSyncer
		encoder zapcore.Encoder
	)

	// 不同的级别的日志输出到不同的地方的判断。决定是否能输出。
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel && lvl >= lc.ParseLevel()
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= lc.ParseLevel()
	})
	encoderConfig := zapcore.EncoderConfig{
		//当存储的格式为JSON的时候这些作为可以key
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//以上字段输出的格式
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if lc.Mode == "console" {
		ws = zapcore.Lock(os.Stdout)
		warnWs = zapcore.Lock(os.Stderr)
		//输出到控制台彩色。
		if !lc.Json {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
	} else {
		normalConfig := &lumberjack.Logger{
			Filename:   lc.FileName,
			MaxSize:    lc.MaxSize,
			MaxAge:     lc.MaxAge,
			MaxBackups: lc.MaxBackup,
			LocalTime:  true,
			Compress:   lc.Compress,
		}
		warnConfig := &lumberjack.Logger{
			Filename:   lc.WarnFileName,
			MaxSize:    lc.MaxSize,
			MaxAge:     lc.MaxAge,
			MaxBackups: lc.MaxBackup,
			LocalTime:  true,
			Compress:   lc.Compress,
		}
		ws = zapcore.AddSync(normalConfig)
		warnWs = zapcore.AddSync(warnConfig)
	}
	if lc.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var core zapcore.Core
	if lc.EnableWarnRedirect {
		highCore := zapcore.NewCore(encoder, warnWs, highPriority)
		lowCore := zapcore.NewCore(encoder, ws, lowPriority)
		core = zapcore.NewTee(highCore, lowCore)
	} else {
		core = zapcore.NewCore(encoder, ws, lc.ParseLevel())
	}
	logger := zap.New(core)
	//是否新增调用者信息
	if lc.AddCaller {
		lc.options = append(lc.options, zap.AddCaller())
	}
	//当错误时是否添加堆栈信息
	if lc.Stacktrace {
		lc.options = append(lc.options, zap.AddStacktrace(zap.PanicLevel))
	}
	return logger.WithOptions(lc.options...)

}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}
