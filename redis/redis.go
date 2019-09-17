package redis

import (
	"time"

	"github.com/go-redis/redis"
)

var (
	redisClient redis.UniversalClient
)

// Configs contains the configuration
// for opening connection to Redis server.
type Configs struct {
	Addresses []string
	Master    string
	Password  string
}

// NewRedisClient creates an instance
// of redis-client, which allow you
// to get data from redis, also set
// new data to it.
func NewRedisClient(cfg Configs) {
	redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		MasterName: cfg.Master,
		Addrs:      cfg.Addresses,
		Password:   cfg.Password,
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
func Set(key string, value interface{}, expiration time.Duration) (string, error) {
	return redisClient.Set(key, value, expiration).Result()
}

// Close closes the connection
// to the redis-server based on the
// current instance in application.
func Close() error {
	return redisClient.Close()
}
