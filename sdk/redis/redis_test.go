package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"testing"
)

func TestInsert(t *testing.T) {
	cli := redis.NewClient(&redis.Options{
		Network:            "",
		Addr:               "192.168.2.200:6379",
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           "local_office_0c1PVg",
		DB:                 0,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolFIFO:           false,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		log.Println(err)
		return
	}
	cli.Set(context.Background(), "admin_kaiPay_config", "", 0)

}
