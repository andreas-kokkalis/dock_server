package config

import (
	"fmt"
	"testing"

	redis "gopkg.in/redis.v5"

	"github.com/stretchr/testify/assert"
)

var validConfigPath = "../../conf/"
var invalidConfigPath = "lalaala"

type confVals struct {
	key, val string
}

// var vals = []confVals{
// 	{"dc.imagerepo.name", "dc"},
// 	{"dc.docker.api.host", "unix:///var/run/docker.sock"},
// 	{"dc.docker.api.version", "1.24"},
// }

// XXX: New Version

func TestNewConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testName := "valid config path"
	c, err := NewConfig(validConfigPath, "local")
	assert.NoError(err, testName)
	assert.NotNil(c)

	testName = "valid config path"
	_, err = NewConfig(invalidConfigPath, "local")
	assert.Error(err, testName)
}

func TestGetPGConnectionString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	expect := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		"localhost",
		"5432",
		"dock",
		"dock",
		"dock",
		"disable",
	)

	testName := "GetPGConnectionString"
	c, err := NewConfig(validConfigPath, "local")
	assert.NoError(err, testName)
	actual := c.GetPGConnectionString()
	assert.Equal(expect, actual, testName)
}

func TestGetRedisConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	expect := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	testName := "GetRedisConfig"
	c, err := NewConfig(validConfigPath, "local")
	assert.NoError(err, testName)
	actual := c.GetRedisConfig()
	assert.Equal(expect, actual, testName)

}

func TestGetDockerConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	expect := map[string]string{
		"host":    "unix:///var/run/docker.sock",
		"version": "1.24",
		"repo":    "dc",
	}

	testName := "GetDockerConfig"
	c, err := NewConfig(validConfigPath, "local")
	assert.NoError(err, testName)
	actual := c.GetDockerConfig()
	assert.Equal(expect, actual, testName)
}

func TestGetAPIPorts(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	expect := 200

	testName := "GetAPIPort"
	c, err := NewConfig(validConfigPath, "local")
	assert.NoError(err, testName)
	actual := c.GetAPIPorts()
	assert.Equal(expect, actual, testName)
}
