package config

import (
	"fmt"
	"testing"

	redis "gopkg.in/redis.v5"

	"github.com/stretchr/testify/assert"
)

var (
	validConfigPath   = "../../conf/"
	invalidConfigPath = "lalaala"

	validEnv = "local"
)

type confVals struct {
	key, val string
}

// var vals = []confVals{
// 	{"dc.imagerepo.name", "dc"},
// 	{"dc.docker.api.host", "unix:///var/run/docker.sock"},
// 	{"dc.docker.api.version", "1.24"},
// }

func newConf(t *testing.T) *Config {
	c, err := NewConfig(validConfigPath, "local")
	assert.NoError(t, err, "Initialize config")
	assert.NotNil(t, c, "Initialize config")
	return c
}

func TestNewConfig(t *testing.T) {
	_ = newConf(t)
	testName := "invalid config path"
	c, err := NewConfig(invalidConfigPath, "local")
	assert.Error(t, err, testName)
	assert.Nil(t, c, testName)
}

func TestGetPGConnectionString(t *testing.T) {
	c := newConf(t)
	expect := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		"localhost",
		"5432",
		"dock",
		"dock",
		"dock",
		"disable",
	)
	actual := c.GetPGConnectionString()
	assert.Equal(t, expect, actual, "GetPGConnectionString")
}

func TestGetRedisConfig(t *testing.T) {
	c := newConf(t)
	expect := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	actual := c.GetRedisConfig()
	assert.Equal(t, expect, actual, "GetRedisConfig")

}

func TestGetDockerConfig(t *testing.T) {
	c := newConf(t)
	expect := map[string]string{
		"host":    "unix:///var/run/docker.sock",
		"version": "1.25",
		"repo":    "andreaskokkalis/dc",
	}
	actual := c.GetDockerConfig()
	assert.Equal(t, expect, actual, "GetDockerConfig")
}

func TestGetAPIPorts(t *testing.T) {
	c := newConf(t)
	expect := 200
	actual := c.GetAPIPorts()
	assert.Equal(t, expect, actual, "GetPortNumbers")
}

func TestGetAPIServerPort(t *testing.T) {
	c := newConf(t)
	expect := ":8080"
	actual := c.GetAPIServerPort()
	assert.Equal(t, expect, actual, "GetAPIServerPort")
}
