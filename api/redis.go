package main

import (
	"context"
	"log"

	"github.com/KonradLinkowski/go-redis/shared"
	"github.com/redis/go-redis/v9"
)

func InitRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: shared.RedisAddr,
	})
}

func PushToRedis(ctx context.Context, rdb *redis.Client, imageName string) error {
	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: shared.RedisStreamName,
		Values: map[string]interface{}{
			shared.RedisStreamValueKey: imageName,
		},
		MaxLen: 100,
		Approx: true,
	}).Result()
	return err
}

func SubscribeToRedis(ctx context.Context, rdb *redis.Client, hub *Hub, handlerFunc func(string)) {
	lastID := "0"

	for {
		streams, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{shared.RedisProcessedStreamName, lastID},
			Block:   0,
		}).Result()

		if err != nil {
			log.Println(err)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				lastID = msg.ID
				imageName := msg.Values[shared.RedisStreamValueKey].(string)
				handlerFunc(imageName)
			}
		}
	}
}
