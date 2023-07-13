package main

import (
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestInsert(t *testing.T) {
	redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}
