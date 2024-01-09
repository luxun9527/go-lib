## zap日志

参考

 https://www.liwenzhou.com/posts/Go/zap/

https://github.com/uber-go/zap/blob/master/example_test.go

https://github.com/flipped-aurora/gin-vue-admin

https://github.com/douyu/jupiter/tree/master/pkg/xlog



本文地址

https://github.com/luxun9527/go-lib/tree/master/utils/zap	 

### 基本概念

决定日志输出的三个参数

1、enc Encoder 决定日志输出的格式(比如日志的时间格式，是否有颜色，日志等级是否大写等)

2、ws WriteSyncer 决定日志输出的位置(输出到标准输出还是文件)

3、 enab LevelEnabler 决定日志输出的日志级别



#### options

一些附加的配置

```go
//是否加上调用者
zap.AddCaller()
//在哪个等级的日志加上堆栈
zap.AddStacktrace(zap.PanicLevel)
//调用链跳过几级，比如你自己封装zap的logger,如果要记住你在哪个地方调用了Info()方法,要加上这个optionAddCallerSkip(1)
func()Info{
    logger.Info()
}
zap.AddCallerSkip(lc.CallerShip)
logger.WithOptions(lc.options...)
```

如果要定制zap的日志，修改对应的配置就好了

```go
// NewCore creates a Core that writes logs to a WriteSyncer.
func NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler) Core {
    return &ioCore{
       LevelEnabler: enab,
       enc:          enc,
       out:          ws,
    }
}

//通过配置创建一个core
core = zapcore.NewCore(encoder, ws, atomicLevel)
//初始化logger
var  logger *zap.Logger  := zap.New(core)
//加上一些option
logger.WithOptions(lc.options...)

l.Debug("this a debug level log", zap.Any("test", "t"))
l.Info("this a info level log", zap.Any("test", "t"))
l.Warn("this a warn level log", zap.Any("test", "t"))
l.Error("this a error level log", zap.Any("test", "t"))
l.Panic("this a panic level log", zap.Any("test", "t"))

2023-11-12-20:41:54	DEBUG	E:/demoproject/go-lib/utils/zap/zap_test.go:33	this a debug level log	{"test": "t"}
2023-11-12-20:41:54	INFO	E:/demoproject/go-lib/utils/zap/zap_test.go:34	this a info level log	{"test": "t"}
2023-11-12-20:41:54	WARN	E:/demoproject/go-lib/utils/zap/zap_test.go:35	this a warn level log	{"test": "t"}
2023-11-12-20:41:54	ERROR	E:/demoproject/go-lib/utils/zap/zap_test.go:36	this a error level log	{"test": "t"}
2023-11-12-20:41:54	PANIC	E:/demoproject/go-lib/utils/zap/zap_test.go:37	this a panic level log	{"test": "t"}
go-lib/utils/zap.TestZap
	E:/demoproject/go-lib/utils/zap/zap_test.go:37
testing.tRunner
	D:/go/src/testing/testing.go:1576
--- FAIL: TestZap (0.01s)
panic: this a panic level log [recovered]
	panic: this a panic level log

```







#### Encoder 

```go
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
if lc.Json {
	encoder = zapcore.NewJSONEncoder(encoderConfig)
} else {
	encoder = zapcore.NewConsoleEncoder(encoderConfig)
}
```

#### WriteSyncer 

```go
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
ws = zapcore.Lock(zapcore.AddSync(normalConfig))
errorWs = zapcore.Lock(zapcore.AddSync(warnConfig))
```

#### LevelEnabler 

```
// LevelEnabler decides whether a given logging level is enabled when logging a
// message.
//
// Enablers are intended to be used to implement deterministic filters;
// concerns like sampling are better implemented as a Core.
//
// Each concrete Level value implements a static LevelEnabler which returns
// true for itself and all higher logging levels. For example WarnLevel.Enabled()
// will return true for WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, and
// FatalLevel, but return false for InfoLevel and DebugLevel.
type LevelEnabler interface {
    Enabled(Level) bool
}
```

```go
// The bundled Config struct only supports the most common configuration
// options. More complex needs, like splitting logs between multiple files
// or writing to non-file outputs, require use of the zapcore package.
//
// In this example, imagine we're both sending our logs to Kafka and writing
// them to the console. We'd like to encode the console output and the Kafka
// topics differently, and we'd also like special treatment for
// high-priority logs.

// First, define our level-handling logic.
highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
    return lvl >= zapcore.ErrorLevel
})
lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
    return lvl < zapcore.ErrorLevel
})
```





## 其他用法

### 异步落盘日志

使用缓存，定时批量落盘日志，能够显著减少磁盘io，提高性能。如果要使用异步落盘日志，记得在程序关闭的使用调用sync()方法

```go
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
```

### 动态调整日志级别

AtomicLevel 是并发安全的，如果有需要在程序运行时调整日志的需求可以使用。

```go
func ExampleAtomicLevel() {
    atom := zap.NewAtomicLevel()

    // To keep the example deterministic, disable timestamps in the output.
    encoderCfg := zap.NewProductionEncoderConfig()
    encoderCfg.TimeKey = ""

    logger := zap.New(zapcore.NewCore(
       zapcore.NewJSONEncoder(encoderCfg),
       zapcore.Lock(os.Stdout),
       atom,
    ))
    defer logger.Sync()

    logger.Info("info logging enabled")

    atom.SetLevel(zap.ErrorLevel)
    logger.Info("info logging disabled")
    // Output:
    // {"level":"info","msg":"info logging enabled"}
}
```

### 日志输出到多个目的地

这个设置能将日志输出到不同的目的地。

当你有将日志同时输出到不同地方可以设置这个，**还可以做到不同的级别的日志输出到不同的地方。**

```go
highCore := zapcore.NewCore(encoder, errorWs, zapcore.ErrorLevel)
lowCore := zapcore.NewCore(encoder, ws, atomicLevel)
core = zapcore.NewTee(highCore, lowCore)
logger := zap.New(core)

```

## 我的设置

分享一个我参考别的开源框架对zap日志的一些设置

```go
package logger

import (
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
	//日志级别 debug info warn panic
	Level string
	//panic时候 是否显示堆栈 panic级别的日志输出堆栈信息。
	Stacktrace bool
	//添加调用者信息
	AddCaller bool
	//调用链，往上多少级 ，在一些中间件，对日志有包装，可以通过这个选项指定。
	CallerShip int
	//输出到哪里标准输出console,还是文件file
	Mode string
	//文件名称加路径
	FileName string
	//error级别的日志输入到不同的地方
	ErrorFileName string
	// 日志文件大小 单位MB 默认500MB
	MaxSize int
	//日志保留天数
	MaxAge int
	//日志最大保留的个数
	MaxBackup int
	//异步日志 日志将先输入到内存到，定时批量落盘。如果设置这个值，要保证在程序退出的时候调用Sync(),在开发阶段不用设置为true。
	Async bool
	//是否 输出json格式的数据，JSON格式相对于console格式，不方便阅读，但是对机器更加友好
	//最佳实践，在开发的时候json为false,mode为console
	Json bool
	//是否日志压缩
	Compress    bool
	options     []zap.Option
	atomicLevel zap.AtomicLevel
}

func (lc *Config) UpdateLevel(level zapcore.Level) {
	lc.atomicLevel.SetLevel(level)
}

func (lc *Config) Build() *zap.Logger {
	if lc.Mode == "file" && lc.FileName == "" {
		log.Printf("file mode, but file name is empty")
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
	if lc.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var (
		core zapcore.Core
	)
	atomicLevel := lc.parseLevel()
	if lc.ErrorFileName != "" && lc.Mode == "file" {
		lowCore := zapcore.NewCore(encoder, ws, atomicLevel)
		c := []zapcore.Core{lowCore}
		if errorWs != nil {
			highCore := zapcore.NewCore(encoder, errorWs, zapcore.ErrorLevel)
			c = append(c, highCore)
		}
		core = zapcore.NewTee(c...)
	} else {
		core = zapcore.NewCore(encoder, ws, atomicLevel)
	}
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
		lc.options = append(lc.options, zap.AddStacktrace(zap.PanicLevel))
	}

	lc.atomicLevel = atomicLevel
	return logger.WithOptions(lc.options...)

}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}

```

如果您要了解zap的基础用法您可以参考一下。如果您深入了解，推荐您直接看zap的example和zap的源码。https://github.com/uber-go/zap/blob/master/example_test.go







