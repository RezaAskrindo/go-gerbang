package database

import (
	"context"

	"go-gerbang/config"

	"github.com/redis/go-redis/v9"
)

var RedisCtx = context.Background()

var RedisDb = redis.NewClient(&redis.Options{
	Addr:         config.Config("REDIS_ADDRES"),
	Network:      config.Config("REDIS_NETWORK"),
	PoolSize:     5, // Small pool limit for a low-spec VPS
	MinIdleConns: 1,
})
