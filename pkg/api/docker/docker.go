package docker

import (
	"os"

	"github.com/docker/docker/client"
)

// DockerCli ...
type DockerCli struct {
	Cli *client.Client
}

// NewAPIClient initializes a new Docker API client.
func NewAPIClient(dockerConfig map[string]string) (*DockerCli, error) {

	docker := &DockerCli{}

	os.Setenv("DOCKER_API_VERSION", dockerConfig["version"])
	os.Setenv("DOCKER_HOST", dockerConfig["host"])

	cli, err := client.NewEnvClient()
	if err != nil {
		docker.Cli = nil
		return docker, err
	}
	docker.Cli = cli
	return docker, nil
}
