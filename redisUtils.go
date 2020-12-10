package main

import "context"
import "github.com/go-redis/redis/v8"

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
	rdb.ZAdd(ctx, "winners", &redis.Z{Score: 0, Member: clientSid})
}
