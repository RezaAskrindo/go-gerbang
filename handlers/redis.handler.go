package handlers

import (
	"context"
	"encoding/json"

	"sika_apigateway/config"

	"github.com/go-redis/redis/v8"
)

// IF PRODUCTION DOMAINESIA USE UNIX
var Cache = redis.NewClient(&redis.Options{
	Addr:    config.Config("REDIS_ADDRES"),
	Network: config.Config("REDIS_NETWORK"),
})

var Ctx = context.Background()

func ToJson(val []byte) (body interface{}) {
	res := body
	err := json.Unmarshal(val, &res)
	if err != nil {
		panic(err)
	}
	return res
}

func ToMarshal(val interface{}) (body []byte) {
	body, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return body
}
