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

func insertTargetSid(sid string) {
	rdb := getClient()
	defer rdb.Close()
	ts := float64(time.Now().Unix())
	rdb.ZAddNX(ctx, "targetsids", &redis.Z{Score: ts, Member: sid})
}

func getTargetSids() []string {
	rdb := getClient()
	defer rdb.Close()
	allSids, err := rdb.ZRevRange(ctx, "targetsids", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	return allSids
}

func clearTargetSids() {
	rdb := getClient()
	defer rdb.Close()
	rdb.Del(ctx, "targetsids")
}

func insertWinner(clientSid string) {
	rdb := getClient()
	defer rdb.Close()
	ts := float64(time.Now().Unix())
	rdb.ZAddNX(ctx, "winners", &redis.Z{Score: ts, Member: clientSid})
}

func getWinners() []interface{} {
	rdb := getClient()
	defer rdb.Close()
	winners, err := rdb.ZRevRangeWithScores(ctx, "winners", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	var output []interface{}
	for _, s := range winners {
		arr := []interface{}{s.Member, s.Score}
		output = append(output, arr)
	}
	return output
}
