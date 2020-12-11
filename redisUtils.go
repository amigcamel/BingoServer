package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func getClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       4,
	})
	return rdb
}

func getTargetSids() []string {
	rdb := getClient()
	allSids, err := rdb.LRange(ctx, "sids", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	return allSids
}

func insertWinner(clientSid string) {
	rdb := getClient()
	ts := float64(time.Now().Unix())
	rdb.ZAddNX(ctx, "winners", &redis.Z{Score: ts, Member: clientSid})
}
