package store

import (
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api/store/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisRepo(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(mock.NewRedis())
	assert.NotNil(t, redisRepo)
}

func TestStripSessionKeyPrefix(t *testing.T) {
	redisRepo := NewRedisRepo(mock.NewRedis())
	expect := "koko"
	actual := redisRepo.StripSessionKeyPrefix("usr:koko")
	assert.Equal(t, expect, actual)
}

func TestGetUserRunKey(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(mock.NewRedis())
	expect := "usr:koko"
	actual := redisRepo.GetUserRunKey("koko")
	assert.Equal(t, expect, actual)
}

func TestCreateAdminKey(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(mock.NewRedis())
	expect := "adm:7ff10abb653dead4186089acbd2b7891"
	actual := redisRepo.CreateAdminKey(1)
	assert.Equal(t, expect, actual)
}

func TestGetAdminSessionRunKey(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(mock.NewRedis())
	expect := "run:1"
	actual := redisRepo.GetAdminSessionRunKey("1")
	assert.Equal(t, expect, actual)
}
