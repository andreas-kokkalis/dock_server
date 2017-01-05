package dc

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/andreas-kokkalis/dock-server/conf"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// CreateContainer creates a container based on
// imageName and reference Tag,
// and returns the containerID
func CreateContainer(imageID, password string) (containerID string, port int, err error) {
	// Set environment variables for shellinabox container
	envVars := []string{"SIAB_PASSWORD=" + password, "SIAB_SUDO=true"}
	// Get the imageTag
	refTag, err := GetTagByID(imageID)
	if err != nil {
		return containerID, port, err
	}

	// --- Container configuration
	// Set container port. This port will be exposed and mapped to a host port
	var natPort nat.Port = "4200/tcp"
	exposedPorts := map[nat.Port]struct{}{natPort: {}}
	// Define configuration required to create a container
	img := conf.GetVal("dc.imagerepo.name") + ":" + refTag
	containerConfig := container.Config{Env: envVars, ExposedPorts: exposedPorts, Image: img}
	// Get a non utilized host port, to avoid collision
	port, err = ContainerPorts.Reserve()
	if err != nil {
		log.Printf("[CreateContainer]: %v", err.Error())
		return "", -1, err
	}
	if port == -1 {
		log.Printf("[CreateContainer]: No ports were available to reserve.\n")
		return "", -1, errors.New("There are no resources available in the system.")
	}

	// --- Host configuration
	// Prepare portBindings containerPort -> Host port. are part of PortMap
	portBindings := []nat.PortBinding{nat.PortBinding{HostPort: strconv.Itoa(port)}}
	// ContainerPorts.PrintUsed() // Debug Logging
	// PortMap is member of container.HostConfig
	portMap := map[nat.Port][]nat.PortBinding{natPort: portBindings}
	hostConfig := container.HostConfig{PortBindings: portMap}

	// Send the request to create the container
	body, err := Cli.ContainerCreate(context.Background(), &containerConfig, &hostConfig, &network.NetworkingConfig{}, "")
	if err != nil {
		ContainerPorts.Remove(port)
		log.Printf("[CreateContainer]: %v", err.Error())
		return "", -1, err
	}

	// Return only the first 12 digits from the sha256 identifier of the container
	return body.ID[:12], port, nil
}
