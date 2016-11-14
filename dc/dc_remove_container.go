package dc

import (
	"context"

	"github.com/docker/docker/api/types"
)

// RemoveContainer force removes a container
func RemoveContainer(containerID string) error {
	options := types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: true, Force: true}
	err := Cli.ContainerRemove(context.Background(), containerID, options)
	if err != nil {
		return err
	}
	return nil
}
