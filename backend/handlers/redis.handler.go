package handlers

import (
	"context"
	"encoding/json"

	"go-gerbang/config"
	"go-gerbang/database"
)

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

func SaveToRedis(Key string, Value interface{}) error {
	cacheDelErr := database.RedisDb.Del(Ctx, Key).Err()
	if cacheDelErr != nil {
		return cacheDelErr
	}

	data := ToMarshal(Value)

	cacheSetErr := database.RedisDb.Set(Ctx, Key, data, config.AuthTimeCache).Err()
	if cacheSetErr != nil {
		return cacheSetErr
	}

	return nil
}
