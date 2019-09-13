package redis

import (
	"time"

	"github.com/go-redis/redis"
)

var (
	redisClient redis.UniversalClient
)

// NewRedisClient creates an instance
// of redis-client, which allow you
// to get data from redis, also set
// new data to it.
func NewRedisClient(master, password string, addrs []string) {
	redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName: master,
		Addrs:      addrs,
		Password:   password,
	})
}

// Get gets value from
// redis-server with a given key.
func Get(key string) (string, error) {
	return redisClient.Get(key).Result()
}

// Set sets the value
// into redis-server based on
// the key.
func Set(key string, value interface{}, expiration time.Duration) {
	redisClient.Set(key, value, expiration)
}

// Close closes the connection
// to the redis-server based on the
// current instance in application.
func Close() error {
	return redisClient.Close()
}