package internal

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ConnectRedis(addr string, password string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
