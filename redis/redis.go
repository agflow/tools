package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/agflow/tools/log"
)

// NewClient returns a verified redis client
func NewClient(address, password string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // use default DB
	})
	if rdb == nil {
		log.Warnf("failed to initialize redis client")
		return rdb
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Warnf("can't ping redis, %v", err)
	}
	return rdb
}
