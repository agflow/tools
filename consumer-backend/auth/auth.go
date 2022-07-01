package auth

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/agflow/tools/agstring"
	"github.com/agflow/tools/consumer-backend/accountmanagement"
	"github.com/agflow/tools/consumer-backend/model"
	"github.com/agflow/tools/log"
)

func getUserFromRedis(rdb *redis.Client, key string) *model.User {
	if rdb == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	val, err := rdb.Get(ctx, "user:"+key).Bytes()
	if err == nil && val != nil {
		user := &model.User{}
		if err := json.Unmarshal(val, user); err == nil {
			return user
		}
		log.Infof("unable to unmarshal user object, %v", err)
	}
	return nil
}

func cacheUser(token string, user *model.User, redisClient *redis.Client) {
	if redisClient == nil {
		return
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Warnf("can't marshal user object, err: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// cache user
	if err := redisClient.Set(ctx, "user:"+token, userBytes, 1*time.Minute).Err(); err != nil {
		log.Warnf("can't set user cache value, err: %v", err)
	}
}

func GetUser(token, cacheOptionsStr string, accountService accountmanagement.API, cacheService *redis.Client) (*model.User, error) {
	cacheOptions := strings.Split(cacheOptionsStr, ",")
	if user := getUserFromRedis(cacheService, token); user != nil && !agstring.ContainsAny(cacheOptions, "no-cache") {
		return user, nil
	}

	user, err := accountService.GetUser(token)
	if err != nil {
		return nil, err
	}

	go cacheUser(token, user, cacheService)

	return user, nil
}
