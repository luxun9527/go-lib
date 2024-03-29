package logger

import (
	"go-lib/utils/zap/report"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"time"
)

var L *zap.Logger

const (
	// _defaultBufferSize specifies the default size used by Buffer.
	_defaultBufferSize = 256 * 1024 // 256 kB

	// _defaultFlushInterval specifies the default flush interval for
	// Buffer.
	_defaultFlushInterval = 30 * time.Second
)

func InitLogger(loggerConfig Config) {
	L = loggerConfig.Build()
}

func (lc *Config) parseLevel() zap.AtomicLevel {
	level, err := zap.ParseAtomicLevel(lc.Level)
	if err != nil {
		log.Panicf("init level failed level %s err %v", lc.Level, err)
	}
	return level
}

type Config struct {
	Name string `json:",default=project_name"`
	//日志级别 debug info warn panic
	Level string `json:",default=debug"`
	//panic时候 是否显示堆栈 panic级别的日志输出堆栈信息。
	Stacktrace bool `json:",default=true"`
	//添加调用者信息
	AddCaller bool `json:",default=true"`
	//调用链，往上多少级 ，在一些中间件，对日志有包装，可以通过这个选项指定。
	CallerShip int `json:",default=3"`
	//输出到哪里标准输出console,还是文件file
	Mode string `json:",default=console"`
	//文件名称加路径
	FileName string `json:",default=console"`
	//error级别的日志输入到不同的地方
	ErrorFileName string `json:",optional"`
	// 日志文件大小 单位MB 默认500MB
	MaxSize int `json:",optional"`
	//日志保留天数 `json:",optional"`
	MaxAge int `json:",optional"`
	//日志最大保留的个数
	MaxBackup int `json:",optional"`
	//异步日志 日志将先输入到内存到，定时批量落盘。如果设置这个值，要保证在程序退出的时候调用Sync(),在开发阶段不用设置为true。
	Async bool `json:",optional"`
	//是否输出json格式
	Json bool `json:",optional"`
	//是否日志压缩
	Compress bool `json:",optional"`
	//是否report
	IsReport     bool             `json:",optional"`
	ReportConfig *report.ImConfig `json:",optional"`
	// 打印到控制台是否带颜色
	Color       bool `json:",default=true"`
	options     []zap.Option
	atomicLevel zap.AtomicLevel
}

func (lc *Config) UpdateLevel(level zapcore.Level) {
	lc.atomicLevel.SetLevel(level)
}

func (lc *Config) Build() *zap.Logger {
	if lc.Mode != "file" && lc.Mode != "console" {
		log.Panicln("mode must be console or file")
	}

	if lc.Mode == "file" && lc.FileName == "" {
		log.Panicln("file mode, but file name is empty")
	}
	var (
		ws      zapcore.WriteSyncer
		errorWs zapcore.WriteSyncer
		encoder zapcore.Encoder
	)
	encoderConfig := zapcore.EncoderConfig{
		//当存储的格式为JSON的时候这些作为可以key
		MessageKey:    "message",
		LevelKey:      "atomicLevel",
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
		if lc.Color {
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
		if lc.ErrorFileName != "" {
			errorConfig := &lumberjack.Logger{
				Filename:   lc.ErrorFileName,
				MaxSize:    lc.MaxSize,
				MaxAge:     lc.MaxAge,
				MaxBackups: lc.MaxBackup,
				LocalTime:  true,
				Compress:   lc.Compress,
			}
			errorWs = zapcore.Lock(zapcore.AddSync(errorConfig))
		}

		ws = zapcore.Lock(zapcore.AddSync(normalConfig))

	}
	if lc.Async {
		ws = &zapcore.BufferedWriteSyncer{
			WS:            ws,
			Size:          _defaultBufferSize,
			FlushInterval: _defaultFlushInterval,
		}
		if errorWs != nil {
			errorWs = &zapcore.BufferedWriteSyncer{
				WS:            errorWs,
				Size:          _defaultBufferSize,
				FlushInterval: _defaultFlushInterval,
			}
		}

	}
	if lc.Color {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var (
		core zapcore.Core
	)
	atomicLevel := lc.parseLevel()

	var c = []zapcore.Core{zapcore.NewCore(encoder, ws, atomicLevel)}
	if errorWs != nil {
		highCore := zapcore.NewCore(encoder, errorWs, zapcore.ErrorLevel)
		c = append(c, highCore)
	}
	if lc.IsReport {

		highCore := zapcore.NewCore(encoder, nil, zapcore.ErrorLevel)
		c = append(c, highCore)
	}

	core = zapcore.NewTee(c...)

	logger := zap.New(core)
	//是否新增调用者信息
	if lc.AddCaller {
		lc.options = append(lc.options, zap.AddCaller())
		if lc.CallerShip != 0 {
			lc.options = append(lc.options, zap.AddCallerSkip(lc.CallerShip))
		}
	}
	//当错误时是否添加堆栈信息
	if lc.Stacktrace {
		lc.options = append(lc.options, zap.AddStacktrace(zap.ErrorLevel))
	}

	lc.atomicLevel = atomicLevel
	return logger.WithOptions(lc.options...)

}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}
