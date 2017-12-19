package store

import (
	"log"
	"time"

	"gopkg.in/redis.v5"
)

// RedisClient interface
type RedisClient interface {
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
	Del(string) (int64, error)
	Exists(string) (bool, error)
}

// Redis is the redis connection
type Redis struct {
	cli *redis.Client
}

// Get implements redis.Get
func (r *Redis) Get(key string) (string, error) {
	return r.cli.Get(key).Result()
}

// Set implements redis.Set
func (r *Redis) Set(key string, value string, expr time.Duration) (string, error) {
	return r.cli.Set(key, value, expr).Result()
}

// Del implements redis.Del
func (r *Redis) Del(key string) (int64, error) {
	return r.cli.Del(key).Result()
}

// Exists implementes redis.Exists
func (r *Redis) Exists(key string) (bool, error) {
	return r.cli.Exists(key).Result()
}

// NewRedisClient ...
func NewRedisClient(options *redis.Options) (*Redis, error) {

	cli := redis.NewClient(options)
	pong, err := cli.Ping().Result()
	if err != nil {
		return &Redis{nil}, err
	}
	log.Printf("Server is responding: %s", pong)
	return &Redis{cli}, nil
}

// Close closes the redis connection
func (r *Redis) Close() error {
	return r.cli.Close()
}
