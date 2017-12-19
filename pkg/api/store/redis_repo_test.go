package store

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/cache/redismock"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisRepo(t *testing.T) {
	redisRepo := NewRedisRepo(redismock.NewRedisMock())
	assert.NotNil(t, redisRepo)
}

func TestStripSessionKeyPrefix(t *testing.T) {
	redisRepo := NewRedisRepo(redismock.NewRedisMock())
	expect := "koko"
	actual := redisRepo.StripSessionKeyPrefix("usr:koko")
	assert.Equal(t, expect, actual)
}

func TestGetUserRunKey(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(redismock.NewRedisMock())
	expect := "usr:koko"
	actual := redisRepo.GetUserRunKey("koko")
	assert.Equal(t, expect, actual)
}

func TestCreateAdminKey(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(redismock.NewRedisMock())
	expect := "adm:7ff10abb653dead4186089acbd2b7891"
	actual := redisRepo.CreateAdminKey(1)
	assert.Equal(t, expect, actual)
}

func TestGetAdminSessionRunKey(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(redismock.NewRedisMock())
	expect := "run:1"
	actual := redisRepo.GetAdminSessionRunKey("1")
	assert.Equal(t, expect, actual)
}

func TestDeleteUserRunConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithDel(1, nil).WithGet("1", nil))
	actual := redisRepo.DeleteUserRunConfig("1")
	assert.NoError(actual)
}

func TestExistsUserRunConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithExists(true, nil))
	actual, err := redisRepo.ExistsUserRunConfig("1")
	assert.NoError(err)
	assert.Equal(true, actual)

	redisRepo = NewRedisRepo(redismock.NewRedisMock().WithExists(false, errors.New("error")))
	actual, err = redisRepo.ExistsUserRunConfig("1")
	assert.Error(err)
	assert.Equal(false, actual)
}

var runCfg = api.RunConfig{
	ContainerID: "asdasx213",
	Port:        "4200",
	Username:    "guest",
	Password:    "password",
	URL:         "test",
}

func TestUserRunConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	v, _ := json.Marshal(runCfg)
	valString := string(v)

	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithGet(valString, nil).WithSet("OK", nil))
	actual, err := redisRepo.GetUserRunConfig("1")
	assert.NoError(err)
	assert.Equal(runCfg, actual)
	err = redisRepo.SetUserRunConfig("1", runCfg)
	assert.NoError(err)

	redisRepo = NewRedisRepo(redismock.NewRedisMock().WithGet("", errors.New("error")).WithSet("Error", nil))
	actual, err = redisRepo.GetUserRunConfig("1")
	assert.Error(err)
	assert.Equal(api.RunConfig{}, actual)

	err = redisRepo.SetUserRunConfig("1", runCfg)
	assert.Error(err)
	assert.Equal(errors.New("Not OK"), err)

	redisRepo = NewRedisRepo(redismock.NewRedisMock().WithSet("Error", errors.New("Not OK")))
	err = redisRepo.SetUserRunConfig("1", runCfg)
	assert.Error(err)
	assert.Equal(errors.New("Not OK"), err)
}

func TestExistsAdminKey(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithExists(true, nil))
	actual, err := redisRepo.ExistsAdminRunConfig("1")
	assert.NoError(err)
	assert.Equal(true, actual)

	redisRepo = NewRedisRepo(redismock.NewRedisMock().WithExists(false, errors.New("error")))
	actual, err = redisRepo.ExistsAdminRunConfig("1")
	assert.Error(err)
	assert.Equal(false, actual)
}

func TestAdminSession(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithExists(true, nil).WithSet("", nil).WithDel(0, nil))
	actual, err := redisRepo.ExistsAdminSession("1")
	assert.NoError(err)
	assert.Equal(true, actual)
	err = redisRepo.SetAdminSession("1")
	assert.NoError(err)
	err = redisRepo.DeleteAdminSession("1")
	assert.NoError(err)

	expectErr := errors.New("error")
	redisRepo = NewRedisRepo(redismock.NewRedisMock().WithExists(false, expectErr).WithSet("", expectErr).WithDel(0, expectErr))
	actual, err = redisRepo.ExistsAdminSession("1")
	assert.Error(err)
	assert.Equal(false, actual)
	err = redisRepo.SetAdminSession("1")
	assert.Error(err)
	assert.Equal(expectErr, err)
	err = redisRepo.DeleteAdminSession("1")
	assert.Error(err)
	assert.Equal(expectErr, err)

}

func TestAdminRunConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	v, _ := json.Marshal(runCfg)
	valString := string(v)

	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithDel(1, nil).WithGet(valString, nil).WithSet("OK", nil))
	actual, err := redisRepo.GetAdminRunConfig("1")
	assert.NoError(err)
	assert.Equal(runCfg, actual)
	err = redisRepo.DeleteAdminRunConfig("1")
	assert.NoError(err)
	err = redisRepo.SetAdminRunConfig("1", runCfg)
	assert.NoError(err)

	expectErr := errors.New("error")
	redisRepo = NewRedisRepo(redismock.NewRedisMock().WithDel(1, expectErr).WithGet(valString, expectErr).WithSet("Not OK", expectErr))
	actual, err = redisRepo.GetAdminRunConfig("1")
	assert.Error(err)
	assert.Equal(api.RunConfig{}, actual)
	assert.Equal(expectErr, err)
	err = redisRepo.DeleteAdminRunConfig("1")
	assert.Error(err)
	assert.Equal(expectErr, err)
	err = redisRepo.SetAdminRunConfig("1", runCfg)
	assert.Error(err)
	assert.Equal(expectErr, err)

	redisRepo = NewRedisRepo(redismock.NewRedisMock().WithSet("Not OK", nil))
	err = redisRepo.SetAdminRunConfig("1", runCfg)
	assert.Error(err)
	assert.Equal(errors.New("Not OK"), err)
}

func TestExistsPort(t *testing.T) {
	t.Parallel()

	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithExists(true, nil))
	exists := redisRepo.ExistsPort("4200")
	assert.Equal(t, true, exists)
}

func TestRemoveIncosistentRedisKeys(t *testing.T) {
	t.Parallel()
	redisRepo := NewRedisRepo(redismock.NewRedisMock().WithDel(1, nil).WithGet("1", nil))
	redisRepo.RemoveIncosistentRedisKeys("1")
}
