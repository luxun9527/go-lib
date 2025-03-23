package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luxun9527/zlog"
	"go.uber.org/zap"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	zlog.DevConfig.Json = true
	zlog.DevConfig.CallerShip = 0
	DevConfig := &zlog.Config{
		Level:      zap.NewAtomicLevelAt(zap.InfoLevel),
		Stacktrace: false,
		AddCaller:  true,
		CallerShip: 1,
		Mode:       zlog.FileMode,
		FileName:   "./test.log",
		//ErrorFileName: "./log/err.log",
		MaxSize:   1,
		MaxAge:    0,
		MaxBackup: 5,
		Json:      true,
	}
	zlog.InitDefaultLogger(DevConfig)
	engine.GET("/ping", func(c *gin.Context) {
		level := c.Query("level")
		switch level {
		case "info":
			zlog.Infof("%s", "ping")
		case "warn":
			zlog.Warnf("%s", "ping")
		case "error":
			zlog.Errorf("%s", "ping")
		}

	})
	engine.Run(":9999")
}
