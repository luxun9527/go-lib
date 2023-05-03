package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type Option func(s *Config)

const (
	// defaultBufferSize sizes the buffer associated with each WriterSync.
	defaultBufferSize = 256 * 1024

	// defaultFlushInterval means the default flush interval
	defaultFlushInterval = 5 * time.Second
)

type Config struct {
	//日志级别
	Level string
	//当错误时，是否显示堆栈
	Stacktrace bool
	//添加调用者信息
	AddCaller bool
	//是否控制台显示，/dev/stdout 与输入文件互斥
	Debug bool
	//文件名称加路径
	FileName string
	//warn 级别的日志输出到不同的地方
	WarnFileName string
	// 日志轮转大小，单位MB，默认500MB
	MaxSize int32
	//日志轮转最大时间，单位day，默认1 day
	MaxAge int32
	//日志轮转个数，默认10
	MaxBackup int32
	//日志轮转周期，默认24 hour
	Interval int32
	//异步日志 异步日志缓存区大小256kb
	Async bool
	//是否 输出json格式的数据，JSON格式相对于console格式，不方便阅读，但是对机器更加友好
	Json    bool
	options []zap.Option
}

func main() {
	c := Config{
		Level:      "info",
		Stacktrace: true,
		AddCaller:  true,
		Debug:      true,
		FileName:   "test",
		MaxSize:    2,
		MaxAge:     1,
		MaxBackup:  1,
		Interval:   2,
		Async:      false,
		Json:       false,
		options:    nil,
	}
	logger := c.Build()
	logger.Debug("这是debug", zap.String("key", "value"))
	logger.Info("这是info", zap.String("key", "value"))
	logger.Warn("这是warn", zap.String("key", "value"))
	logger.Error("这是error", zap.String("key", "value"))
	logger.Panic("这是panic", zap.String("key", "value"))
}
func (c *Config) GetLevel() zap.AtomicLevel {
	lv := zap.NewAtomicLevel()
	if err := lv.UnmarshalText([]byte(c.Level)); err != nil {
		panic(err)
	}
	return lv
}

func (c *Config) Build() *zap.Logger {
	var (
		ws      zapcore.WriteSyncer
		encoder zapcore.Encoder
		logger  *zap.Logger
	)
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
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	//如果是debug模式在控制台输出os.stdout
	if c.Debug {
		ws = zapcore.Lock(os.Stdout)
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		//是否异步
		fileConfig := &lumberjack.Logger{
			Filename:   c.FileName,
			MaxSize:    int(c.MaxSize),
			MaxAge:     int(c.MaxAge),
			MaxBackups: int(c.MaxBackup),
			LocalTime:  true,
			Compress:   true,
		}
		if c.Async {
			ws = &zapcore.BufferedWriteSyncer{
				WS: zapcore.AddSync(fileConfig), FlushInterval: defaultFlushInterval, Size: defaultBufferSize}
		} else {
			ws = zapcore.AddSync(fileConfig)
		}
	}

	//是否json格式输出
	if c.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(encoder, ws, c.GetLevel())

	logger = zap.New(core)

	//是否新增调用者信息
	if c.AddCaller {
		c.options = append(c.options, zap.AddCaller())
	}
	//当错误时是否添加堆栈信息
	if c.Stacktrace {
		c.options = append(c.options, zap.AddStacktrace(zap.ErrorLevel))
	}
	logger.WithOptions(c.options...)
	return logger

}
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 - 15:04:05"))
}
