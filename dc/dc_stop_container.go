package dc

import (
	"context"
	"errors"
	"time"

	"github.com/docker/docker/api/types"
)

// StopContainer stops a running container
func StopContainer(containerID string) error {
	t := time.Duration(10)

	// if container is not running return that it was already stopped

	isRunning, err := CheckState(containerID, types.ContainerState{Running: true})
	if err != nil {
		return err
	}
	if !isRunning {
		return nil
	}
	err = Cli.ContainerStop(context.Background(), containerID, &t)
	if err != nil {
		return err
	}

	var isStopped bool
	isStopped, err = CheckState(containerID, types.ContainerState{Running: false})
	if err != nil {
		return err
	}
	if !isStopped == true {
		return errors.New("Container was not stopped")
	}
	return nil
}
