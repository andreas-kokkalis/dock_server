package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAPIClient(t *testing.T) {
	cli, err := NewAPIClient(map[string]string{
		"version": "1.25",
		"host":    "unix:///var/run/docker.sock",
	})
	assert.NoError(t, err)
	assert.NotNil(t, cli)
	cli2, err := NewAPIClient(map[string]string{
		"version": "x",
		"host":    "local",
	})
	assert.Error(t, err)
	assert.Nil(t, cli2)
}
