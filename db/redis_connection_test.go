package db

import "testing"

func TestInitRedisClient(t *testing.T) {
	InitRedisClient()

	pong, err := RCli.Ping().Result()
	if err != nil {
		t.Errorf("Unable to ping the redis server. Response: %s", pong)
	}
}
