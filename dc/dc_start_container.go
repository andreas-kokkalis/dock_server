package dc

import (
	"context"

	"github.com/docker/docker/api/types"
)

// StartContainer sends a request to start a container
func StartContainer(containerID string) error {

	// Start container
	err := Cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		// No need to validate if port number is -1, since error is returned in that case
		return err
	}

	// Check if container is running
	var isRunning bool
	isRunning, err = CheckState(containerID, types.ContainerState{Running: true})
	if err != nil {
		return err
	}
	if isRunning == false {
		//XXX: Need to extend this to wait for a container to run
	}

	return nil
}
