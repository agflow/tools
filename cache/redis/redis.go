package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/agflow/tools/log"
)

// Client is a wrapper of a redis client
type Client struct {
	Redis *redis.Client
}

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

// New returns a new redis client
func New(address, password string) *Client {
	redisClient := NewClient(address, password)
	return &Client{Redis: redisClient}
}

// Get gets the value stored on `key`
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	return c.Redis.Get(ctx, key).Bytes()
}

// Set sets `value` on `key` with a `timeout`
func (c *Client) Set(
	ctx context.Context,
	key string,
	value interface{},
	timeout time.Duration,
) error {
	return c.Redis.Set(ctx, key, value, timeout).Err()
}
