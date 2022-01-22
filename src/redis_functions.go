package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var domainName string = "https://fich.is/"

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
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
