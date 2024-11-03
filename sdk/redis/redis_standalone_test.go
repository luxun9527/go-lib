package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
)

func TestStandalone(t *testing.T) {
	ctx := context.Background()

	// 创建 Redis 集群客户端
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.2.159:6379", // Redis 节点的地址
		Password: "123456",             // 如果设置了密码，可以在这里配置
		DB:       0,                    // 默认 DB 选择，0 是默认数据库
	})
	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis cluster: %v", err)
	}
	log.Printf("Connected to Redis cluster: %s", pong)

	// 使用示例：设置和获取一个 key
	err = rdb.Set(ctx, "example_key", "Hello, Redis!", 0).Err()
	if err != nil {
		log.Fatalf("Error setting key: %v", err)
	}

	val, err := rdb.Get(ctx, "example_key").Result()
	if err != nil {
		log.Fatalf("Error getting key: %v", err)
	}

	log.Printf("Value of 'example_key': %s", val)

	// 关闭客户端
	if err := rdb.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	}
}
