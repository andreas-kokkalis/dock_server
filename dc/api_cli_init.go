package dc

import (
	"os"

	"github.com/docker/docker/client"
)

// DefaultAPIVersion is the version of the Remote Docker API
const DefaultAPIVersion string = "1.24"

// DefaultDockerHost is the connection string for Docker host
const DefaultDockerHost string = "unix:///var/run/docker.sock"

// Cli the global docker api client
var Cli *client.Client

// ClientInit initializes a new client API variable that is globally
// accessible when invoking this package.
func ClientInit(APIVersion string, DockerHost string) {

	if APIVersion == "" {
		os.Setenv("DOCKER_API_VERSION", DefaultAPIVersion)
	} else {
		os.Setenv("DOCKER_API_VERSION", APIVersion)
	}
	if DockerHost == "" {
		os.Setenv("DOCKER_HOST", DefaultDockerHost)
	} else {
		os.Setenv("DOCKER_HOST", DefaultDockerHost)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	Cli = cli
}
