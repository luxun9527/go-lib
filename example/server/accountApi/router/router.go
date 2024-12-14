package router

import (
	"github.com/gin-gonic/gin"
	"github.com/luxun9527/zlog"
	"go-lib/example/pkg/xgin"
	"go.uber.org/zap/zapcore"
)

func InitRouter(engine *gin.Engine) {
	engine.Use(gin.RecoveryWithWriter(zlog.NewWriter(zlog.DefaultLogger, zapcore.ErrorLevel), func(c *gin.Context, err any) {
		zlog.Errorf("recovery from panic err %v", err)
		xgin.FailWithLang(c)
	}))
	engine.Use(zlog.GetGinLogger())
	group := engine.Group("/api/v1")
	initAccountRouter(group)
}
