package main

import (
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestInsert(t *testing.T) {
	redis.Client{}.HGetAll().Scan()

}
