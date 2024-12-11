package database

import (
	"github.com/redis/go-redis/v9"
	"go-gerbang/config"
)

// var RedisAddrs = config.Config("REDIS_ADDRES")
// var RedisAddrsFull = config.Config("REDIS_ADDRESS_FULL")

var RedisDb = redis.NewClient(&redis.Options{
	Addr:    config.Config("REDIS_ADDRES"),
	Network: config.Config("REDIS_NETWORK"),
})
