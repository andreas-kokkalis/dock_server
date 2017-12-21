package docker

import (
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/stretchr/testify/assert"
)

var validConfigDir = "../../../conf"

func TestNewAPIClinet(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	c, _ := config.NewConfig(validConfigDir, "local")
	cli, err := NewAPIClient(c.GetDockerConfig())
	assert.NotNil(cli)
	assert.NoError(err)

	cli2, err := NewAPIClient(map[string]string{"repo": "1", "version": "x", "host": "local"})
	assert.Error(err)
	assert.Nil(cli2.Cli)
}
