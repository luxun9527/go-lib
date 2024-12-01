package xgorm

import (
	"github.com/luxun9527/zlog"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

type GormConf struct {
	Ip           string       `json:"ip"`
	Port         int32        `json:"port"`
	Username     string       `json:"username"`
	Password     string       `json:"password"`
	DbName       string       `json:"dbname"`
	MaxIdleConns int          `json:"maxIdleConns"`
	MaxOpenConns int          `json:"maxOpenConns"`
	Logger       *zlog.Config `json:"logger,optional"`
}

func (gc *GormConf) dns() string {
	return gc.Username + ":" + gc.Password + "@tcp(" + gc.Ip + ":" + cast.ToString(gc.Port) + ")/" + gc.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
}
func (gc *GormConf) MustNewGormClient() *gorm.DB {
	if db, err := gorm.Open(mysql.Open(gc.dns()), gc.gormConfig()); err != nil {
		zlog.Panicf("gorm.Open error: %v", err)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(gc.MaxIdleConns)
		sqlDB.SetMaxOpenConns(gc.MaxOpenConns)
		return db
	}
}

func (gc *GormConf) gormConfig() *gorm.Config {
	if gc.Logger == nil {
		zlog.Panic("gorm logger config is nil")
	}
	config := &gorm.Config{}
	config.SkipDefaultTransaction = true

	if gc.Logger.Level.Level() == zapcore.InfoLevel {
		// 如果是info级别，则设置为debug gorm info級別使用的是debug打印的
		gc.Logger.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	gc.Logger.CallerShip = 3
	l := gc.Logger.Build()
	gl := zapgorm2.New(l)
	gl.IgnoreRecordNotFoundError = true
	gl.LogLevel = gormlogger.Info
	config.Logger = gl

	return config
}

type E struct {
	zap.Logger
}
