package main

import (
	"HiChat/initialize"
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

var ctx = context.Background()

func TestRedisConn(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
	})
	defer client.Close()
	initialize.RedisDB.Get(ctx, "name")
}
