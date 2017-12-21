package docker

import (
	"os"

	"github.com/docker/docker/client"
)

// APIClient ...
type APIClient struct {
	Cli *client.Client
}

// NewAPIClient initializes a new Docker API client.
func NewAPIClient(dockerConfig map[string]string) (*APIClient, error) {

	docker := &APIClient{}

	_ = os.Setenv("DOCKER_API_VERSION", dockerConfig["version"])
	_ = os.Setenv("DOCKER_HOST", dockerConfig["host"])

	cli, err := client.NewEnvClient()
	if err != nil {
		docker.Cli = nil
		return docker, err
	}
	docker.Cli = cli
	return docker, nil
}
