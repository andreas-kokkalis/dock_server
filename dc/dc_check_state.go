package dc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
)

// CheckState checks if container ha
func CheckState(containerID string, state string) (bool, error) {

	var inspect types.ContainerJSON
	var err error

	for i := 0; i < 50; i++ {
		time.Sleep(time.Millisecond)
		inspect, err = Cli.ContainerInspect(context.Background(), containerID)
		if err != nil {
			return false, errors.New("Error inspecting state of container: " + containerID)
		}
		fmt.Printf("Container: %s, Status: %s\n", containerID, inspect.State.Status)
		if inspect.State.Status == state {
			return true, nil
		}
	}
	// After X miliseconds container has not started
	return false, nil
}
