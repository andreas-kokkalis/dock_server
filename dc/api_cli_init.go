package dc

import (
	"os"

	"github.com/andreas-kokkalis/dock-server/conf"
	"github.com/docker/docker/client"
)

// Cli the global docker api client
var Cli *client.Client

// ClientInit initializes a new client API variable that is globally
// accessible when invoking this package.
func ClientInit(apiVersion string, dockerHost string) {

	if apiVersion == "" {
		os.Setenv("DOCKER_API_VERSION", conf.GetVal("dc.docker.api.host"))
	} else {
		os.Setenv("DOCKER_API_VERSION", apiVersion)
	}
	if dockerHost == "" {
		os.Setenv("DOCKER_HOST", conf.GetVal("dc.docker.api.version"))
	} else {
		os.Setenv("DOCKER_HOST", dockerHost)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	Cli = cli
}
