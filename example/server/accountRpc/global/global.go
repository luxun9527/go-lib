package global

import (
	"github.com/redis/go-redis/v9"
	"go-lib/example/server/accountRpc/config"
	accountQuery "go-lib/example/server/accountRpc/dao/account/query"
)

var (
	AccountDB *accountQuery.Query
	Config    *config.Config
	RedisCli  *redis.Client
)
