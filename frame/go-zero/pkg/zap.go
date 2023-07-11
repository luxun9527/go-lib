package pkg

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"time"
)

const (
	// _defaultBufferSize specifies the default size used by Buffer.
	_defaultBufferSize = 256 * 1024 // 256 kB

	// _defaultFlushInterval specifies the default flush interval for
	// Buffer.
	_defaultFlushInterval = 30 * time.Second
)

func (lc *ZapLogger) ParseLevel() zapcore.Level {
	level, err := zapcore.ParseLevel(lc.Level)
	if err != nil {
		log.Panicf("init level failed level %s err %v", lc.Level, err)
	}
	return level
}

type ZapLogger struct {
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
	//error 级别的日志输出到不同的地方
	ErrorFileName string `json:",optional"`
	// 日志轮转大小，单位MB，默认500MB
	MaxSize int
	//日志最多保留天数
	MaxAge int
	//日志最大保留的个数
	MaxBackup int
	//异步日志 待实现
	Async bool `json:",optional"`
	//是否 输出json格式的数据，JSON格式相对于console格式，不方便阅读，但是对机器更加友好
	//最佳实践，在开发的时候json为false,mode为console,测试部署阶段，mode为file,json为TRUE
	Json bool
	//是否日志压缩
	Compress bool
	options  []zap.Option
}

func (lc *ZapLogger) Build() *zap.Logger {

	var (
		ws      zapcore.WriteSyncer
		errorWs zapcore.WriteSyncer
		encoder zapcore.Encoder
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
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if lc.Mode == "console" {
		ws = zapcore.Lock(os.Stdout)
		errorWs = zapcore.Lock(os.Stderr)
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
			Filename:   lc.ErrorFileName,
			MaxSize:    lc.MaxSize,
			MaxAge:     lc.MaxAge,
			MaxBackups: lc.MaxBackup,
			LocalTime:  true,
			Compress:   lc.Compress,
		}
		ws = zapcore.AddSync(normalConfig)
		errorWs = zapcore.AddSync(warnConfig)
	}
	if lc.Async {
		ws = &zapcore.BufferedWriteSyncer{
			WS:            ws,
			Size:          _defaultBufferSize,
			FlushInterval: _defaultFlushInterval,
		}
		errorWs = &zapcore.BufferedWriteSyncer{
			WS:            errorWs,
			Size:          _defaultBufferSize,
			FlushInterval: _defaultFlushInterval,
		}
	}
	if lc.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var core zapcore.Core
	if lc.ErrorFileName != "" && lc.Mode == "file" {
		c1 := zapcore.NewCore(encoder, errorWs, zap.ErrorLevel)
		c2 := zapcore.NewCore(encoder, ws, lc.ParseLevel())
		core = zapcore.NewTee(c1, c2)
	} else {
		core = zapcore.NewCore(encoder, ws, lc.ParseLevel())
	}
	logger := zap.New(core)
	//是否新增调用者信息
	if lc.AddCaller {
		//当错误时是否添加堆栈信息
		lc.options = append(lc.options, zap.AddCaller())
	}
	if lc.Stacktrace {
		lc.options = append(lc.options, zap.AddStacktrace(zap.PanicLevel))
	}

	return logger.WithOptions(lc.options...)

}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}

const callerSkipOffset = 3

type ZapWriter struct {
	logger *zap.Logger
}

func NewZapWriter(conf *ZapLogger) logx.Writer {
	logger := conf.Build()
	return &ZapWriter{
		logger: logger,
	}
}

func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

func (w *ZapWriter) Close() error {
	return w.logger.Sync()
}

func (w *ZapWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Panic(fmt.Sprint(v))
}

func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toZapFields(fields...)...)
}

func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}
func Error(err error) logx.LogField {
	return logx.Field("err", err)
}
func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
