package redis

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/redis.v5"
)

//go:generate moq -out redismock/redismock.go -pkg redismock . Redis

// Redis interface models the basic interactions with Redis required in this project
type Redis interface {
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
	Del(string) (int64, error)
	Exists(string) (bool, error)
	Close() error
}

// Client is the redis connection
type Client struct {
	cli *redis.Client
}

// NewClient ...
func NewClient(options *redis.Options) (Redis, error) {

	cli := redis.NewClient(options)
	_, err := cli.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "Redis server is not responding")
	}
	return &Client{cli}, nil
}

// Get implements redis.Get
func (r *Client) Get(key string) (string, error) {
	return r.cli.Get(key).Result()
}

// Set implements redis.Set
func (r *Client) Set(key string, value string, expr time.Duration) (string, error) {
	return r.cli.Set(key, value, expr).Result()
}

// Del implements redis.Del
func (r *Client) Del(key string) (int64, error) {
	return r.cli.Del(key).Result()
}

// Exists implementes redis.Exists
func (r *Client) Exists(key string) (bool, error) {
	return r.cli.Exists(key).Result()
}

// Close closes the redis connection
func (r *Client) Close() error {
	return r.cli.Close()
}
