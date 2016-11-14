package dc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/andreas-kokkalis/dock-server/srv"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// CreateContainer creates a container based on
// imageName and reference Tag,
// and returns the containerID
func CreateContainer(image, refTag, username, password string) (string, int, error) {

	envVars := []string{"SIAB_PASSWORD=" + password, "SIAB_USER=" + username, "SIAB_SUDO=true"}
	var natPort nat.Port = "4200/tcp"

	// ExposedPorts used in container.Config
	exposedPorts := map[nat.Port]struct{}{natPort: {}}
	containerConfig := container.Config{Env: envVars, ExposedPorts: exposedPorts, Image: image + ":" + refTag}

	// Get a non utilized host port, to avoid collision
	port, err := srv.GetFreeResource(srv.PortResources)
	if err != nil {
		return "", -1, err
	}

	// portBindings are part of PortMap
	portBindings := []nat.PortBinding{nat.PortBinding{HostPort: strconv.Itoa(port)}}
	srv.PrintUsed(srv.PortResources)

	// PortMap is member of container.HostConfig
	portMap := map[nat.Port][]nat.PortBinding{natPort: portBindings}
	hostConfig := container.HostConfig{PortBindings: portMap}

	// Send the request to create the container
	body, err := Cli.ContainerCreate(context.Background(), &containerConfig, &hostConfig, &network.NetworkingConfig{}, "")
	if err != nil {
		srv.FreeResource(srv.PortResources, port)
		return "", -1, err
	}
	fmt.Println(body)
	return body.ID, port, nil
}
