package config

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
		rdb *redis.Client
		sessionExp   = time.Hour * 24
		redisCtx     = context.Background()
	)

func RedisInit() error {
	dbRedis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       3,  // database 3
	})

	rdb = dbRedis
	return nil
}

func RedisConnect() *redis.Client {
	return rdb
}

func GetsessionExp() time.Duration {
	return sessionExp
}

func GetRedisCtx() context.Context {
	return redisCtx
}