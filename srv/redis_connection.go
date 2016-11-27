package srv

import (
	"fmt"

	"gopkg.in/redis.v5"
)

// RCli is the redis client connection
var RCli *redis.Client

// InitRedisClient initializes the redis connection
func InitRedisClient() {
	RCli = redis.NewClient(&redis.Options{
		Addr:     "179.16.238.10:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := RCli.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server is responding: %s", pong)
}
