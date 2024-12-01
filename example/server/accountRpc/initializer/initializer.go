package initializer

import (
	"github.com/luxun9527/zlog"
	"go-lib/example/pkg/xetcd"
	"go-lib/example/pkg/xjwt"
	"go-lib/example/server/accountRpc/config"
	"go-lib/example/server/accountRpc/dao/account/query"
	"go-lib/example/server/accountRpc/global"
)

func Init(configPath string) {
	//初始化配置
	global.Config = config.InitConfig(configPath)

	//初始化日志
	zlog.InitDefaultLogger(&global.Config.Logger)

	//初始化gormdb连接
	db := global.Config.GormConf.MustNewGormClient()
	global.AccountDB = query.Use(db)

	//初始化etcd注册
	xetcd.Register(global.Config.Server.RegisterConf)

	//初始化jwt
	xjwt.DefaultJwtConf = &global.Config.JwtConf
	//初始化redis
	global.RedisCli = global.Config.RedisConf.MustBuildNode()
}
