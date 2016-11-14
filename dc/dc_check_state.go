package dc

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
)

// CheckState checks if container ha
func CheckState(containerID string, state types.ContainerState) (bool, error) {

	var inspect types.ContainerJSON
	var err error

	inspect, err = Cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return false, errors.New("Error inspecting state of container: " + containerID)
	}
	fmt.Println(inspect)

	if inspect.State.Status == state.Status {
		fmt.Printf("Container: %s, Status: %s\n", containerID, inspect.State.Status)
		return true, nil
	}

	return false, nil
}
