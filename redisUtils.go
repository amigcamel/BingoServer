package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func getClient(db ...string) *redis.Client {
	var dbStr string
	if len(db) == 1 {
		dbStr = db[0]
	} else {
		dbStr = os.Getenv("BINGO_REDIS_DB")
	}
	dbInt, err := strconv.Atoi(dbStr)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("BINGO_REDIS_ADDR"),
		Password: "",
		DB:       dbInt,
	})
	return rdb
}

func insertTargetSid(sid string) {
	rdb := getClient()
	defer rdb.Close()
	ts := float64(time.Now().UnixNano() / 1e6)
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

func clearWinners() {
	rdb := getClient()
	defer rdb.Close()
	rdb.Del(ctx, "winners")
}

func insertWinner(clientSid string) int64 {
	rdb := getClient()
	defer rdb.Close()
	ts := float64(time.Now().UnixNano() / 1e6)
	rdb.ZAddNX(ctx, "winners", &redis.Z{Score: ts, Member: clientSid})
	rank, err := rdb.ZRank(ctx, "winners", clientSid).Result()
	if err != nil {
		panic(err)
	}
	return rank
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
