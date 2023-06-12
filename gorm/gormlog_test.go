package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm/model"
	"log"
	"moul.io/zapgorm2"
	"os"
	"testing"
	"time"
)

func TestGormLog(t *testing.T) {
	c := ZapConfig{
		Level:              "debug",
		Stacktrace:         true,
		AddCaller:          true,
		Mode:               "file",
		FileName:           "test.log",
		EnableWarnRedirect: true,
		WarnFileName:       "warn.log",
		MaxSize:            12,
		MaxAge:             1,
		MaxBackup:          4,
		Async:              false,
		Json:               true,
		Compress:           true,
		options:            nil,
	}
	l := c.Build()
	l = l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	logger := zapgorm2.New(l)
	logger.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold:             time.Second, // Slow SQL threshold
	//		LogLevel:                  logger.Info, // Log level
	//		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
	//		Colorful:                  true,        // Disable color
	//
	//	},
	//)
	dsn := "root:123456@tcp(192.168.2.99:3306)/gormtest?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger,
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlBb, err := db.DB()
	sqlBb.SetMaxIdleConns(100)
	sqlBb.SetMaxOpenConns(1000)
	db = db
	user := &model.User{Username: "zhangsan", Fav: "篮球"}

	//INSERT INTO `user` (`username`,`fav`) VALUES ('zhangsan','')
	//	INSERT INTO `user` (`fav`) VALUES ('篮球')
	if err := db.Select("fav").Create(user).Error; err != nil {
		log.Println("create err", err)
	}
	//INSERT INTO `user` (`username`) VALUES ('zhangsan')
	if err := db.Omit("fav").Create(user).Error; err != nil {
		log.Println("create err", err)
	}
}

func main() {
	c := ZapConfig{
		Level:              "",
		Stacktrace:         false,
		AddCaller:          false,
		Mode:               "",
		FileName:           "",
		EnableWarnRedirect: false,
		WarnFileName:       "",
		MaxSize:            0,
		MaxAge:             0,
		MaxBackup:          0,
		Async:              false,
		Json:               false,
		Compress:           false,
		options:            nil,
	}
	logger := c.Build()
	logger.Debug("这是debug", zap.String("key", "value"))
	logger.Info("这是info", zap.String("key", "value"))
	logger.Warn("这是warn", zap.String("key", "value"))
	logger.Error("这是error", zap.String("key", "value"))
	logger.Panic("这是panic", zap.String("key", "value"))
}
func (lc *ZapConfig) ParseLevel() zapcore.Level {
	level, err := zapcore.ParseLevel(lc.Level)
	if err != nil {
		log.Panicf("init level failed level %s err %v", lc.Level, err)
	}
	return level
}

type ZapConfig struct {
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

func (lc *ZapConfig) Build() *zap.Logger {

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
