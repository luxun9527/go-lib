package main

import (
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
		Mode:                "console",
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
	logx.Severe("test zhangsan")
	logx.ErrorStack("zhangsan")
	logx.Stat("zhangsan")

	logx.Slow("test")
}
func TestRedis(t *testing.T) {
	conf := redis.RedisConf{
		Host:        "192.168.2.159:6379",
		Type:        "node",
		Pass:        "",
		Tls:         false,
		NonBlock:    false,
		PingTimeout: time.Hour,
	}
	cli := redis.MustNewRedis(conf)
	lock := redis.NewRedisLock(cli, "lock1")
	lock.SetExpire(10)
	ok, err := lock.Acquire()
	if err != nil {
		log.Println(err)
		return
	}
	if ok {
		log.Println("lock success")
	}
	ok1, err := lock.Release()
	if ok1 {
		log.Println("release success")
	}
}
func TestRedis2(t *testing.T) {
	Sum(1)
}
func Sum(n int) {
	defer log.Println(n)
	n += 122
}
