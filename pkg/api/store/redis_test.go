package store

import (
	"net"
	"testing"

	redis "gopkg.in/redis.v5"

	"github.com/stretchr/testify/assert"
)

func TestInitRedisClient(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testName := "NewRedisClient"
	redis, err := NewRedisClient(
		&redis.Options{
			Addr: ":1234",
			Dialer: func() (net.Conn, error) {
				return net.Dial("tcp", ":6379")
			},
			DB: 0,
		})
	assert.Error(err, testName)
	assert.Nil(redis.cli, testName)
}
