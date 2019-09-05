package redis

import (
	"github.com/go-redis/redis"
	"time"
)

var (
	redisClient redis.UniversalClient
)

func NewRedis(master, password string, addrs []string) {
	redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName:         master,
		Addrs:              addrs,
		Password:           password,
	})
}

func Get(key string) (string, error) {
	return redisClient.Get(key).Result()
}

func Set(key string, value interface {}, expiration time.Duration) {
	redisClient.Set(key, value, expiration)
}