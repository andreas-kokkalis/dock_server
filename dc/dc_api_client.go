package dc

import (
	"os"

	"github.com/docker/docker/client"
)

// Cli the global docker api client
var Cli *client.Client

// APIClientInit initializes a new Docker API client.
func APIClientInit(apiVersion string, dockerHost string) {

	os.Setenv("DOCKER_API_VERSION", apiVersion)
	os.Setenv("DOCKER_HOST", dockerHost)

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	Cli = cli
}
