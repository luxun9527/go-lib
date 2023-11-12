## zap

参考

 https://www.liwenzhou.com/posts/Go/zap/

https://github.com/uber-go/zap/blob/master/example_test.go

https://github.com/flipped-aurora/gin-vue-admin

https://github.com/douyu/jupiter/tree/master/pkg/xlog



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

## 代码

这个是我参考一些开源框架，自己对zap的一些封装[config.go](https://github.com/luxun9527/go-lib/blob/master/utils/zap/config.go)，如果您要了解zap的基础用法您可以参考一下。如果您深入了解，推荐您直接看zap的example和zap的源码。https://github.com/uber-go/zap/blob/master/example_test.go





