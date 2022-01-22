package main

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var domainName string = "https://fich.is/"

var redisHost string = os.Getenv("FICHIS_REDIS_HOST")
var redisPort string = os.Getenv("FICHIS_REDIS_PORT")

var redisAddress string = redisHost + ":" + redisPort

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     redisAddress,
	Password: "",
	DB:       0,
})

func addLink(key string, value string) (link string, err error) {

	link = domainName + key
	err = rdb.Set(ctx, key, value, 0).Err()
	return
}

func deleteLink(key string) (err error) {
	err = rdb.Del(ctx, key).Err()
	return
}

func getLink(key string) (link string, err error) {
	operation := rdb.Get(ctx, key)
	link, err = operation.Val(), operation.Err()
	return
}
