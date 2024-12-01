package xredis

import (
	"context"
	"errors"
	"github.com/luxun9527/zlog"
	"github.com/redis/go-redis/v9"
	"time"
)

// RedisConf v9使用redis7
type RedisConf struct {
	Host         []string
	Password     string
	PoolSize     int
	MinIdleConns int
	PingTimeOut  int64
}

func (rc *RedisConf) MustBuildNode() *redis.Client {
	cli, err := rc.BuildNode()
	if err != nil {
		zlog.Panicf("init redis client failed , error: %v", err)
	}
	return cli
}
func (rc *RedisConf) BuildNode() (*redis.Client, error) {
	if len(rc.Host) == 0 {
		return nil, errors.New("redis host is empty")
	}
	cli := redis.NewClient(&redis.Options{
		Addr:         rc.Host[0],
		Password:     rc.Password, // no password set
		PoolSize:     rc.PoolSize,
		MinIdleConns: rc.MinIdleConns,
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*time.Duration(rc.PingTimeOut))
	defer cancelFunc()
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return cli, nil
}
func (rc *RedisConf) BuildCluster() (*redis.ClusterClient, error) {
	if len(rc.Host) == 0 {
		return nil, errors.New("redis host is empty")
	}
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        rc.Host,
		Password:     rc.Password,
		PoolSize:     rc.PoolSize,
		MinIdleConns: rc.MinIdleConns,
		MaxIdleConns: rc.MinIdleConns,
	})
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*time.Duration(rc.PingTimeOut))
	defer cancelFunc()
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return cli, nil
}
