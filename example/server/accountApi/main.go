package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luxun9527/zlog"
	"github.com/spf13/cast"
	"go-lib/example/server/accountApi/global"
	"go-lib/example/server/accountApi/initializer"
	"go-lib/example/server/accountApi/router"
)

var (
	path = flag.String("f", "example/server/accountApi/conf/config.yaml", "config file path")
)

func main() {
	flag.Parse()
	initializer.Init(*path)
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	router.InitRouter(e)
	zlog.Infof("account api server success on %v", global.Config.Server.Port)
	if err := e.Run(fmt.Sprintf("0.0.0.0:" + cast.ToString(global.Config.Server.Port))); err != nil {
		zlog.Panicf("startup service failed, err:%v", err)
	}
}
