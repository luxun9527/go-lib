package main

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"log"
	"testing"
	"time"
)

func TestLogx(t *testing.T) {
	conf := logx.LogConf{
		ServiceName:         "test-logx",
		Mode:                "file",
		Encoding:            "json",
		TimeFormat:          "",
		Path:                "./logs",
		Level:               "info",
		MaxContentLength:    1,
		Compress:            false,
		Stat:                true,
		KeepDays:            1,
		StackCooldownMillis: 0,
		MaxBackups:          1,
		MaxSize:             1,
		Rotation:            "daily",
	}
	logc.MustSetup(conf)
	logx.Debug("zhangsan")
	logx.Info("zhangsan")
	logx.Error("zhangsan")
	logx.ErrorStack("zhangsan")
	logx.Stat("zhangsan")

	logx.Slow("test")
}
func TestRedis(t *testing.T) {
	conf := redis.RedisConf{
		Host:        "192.168.179.99:6379",
		Type:        "node",
		Pass:        "",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Hour,
	}
	cli := redis.MustNewRedis(conf)
	val, err := cli.Hget("test", "zha1ngsan")
	if err != nil {
		ok := errors.Is(err, redis.Nil)
		log.Println(ok)
		return
	}

	log.Println(val)
}
