package dc

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// GetContainers returns the list of containers. Use
// type.ContainerListOptions to filter for state such as
// status=(created,	restarting, running, paused, exited, dead)
func GetContainers(status string) ([]Ctn, error) {
	var containerList []Ctn

	// If containers are filtered by status, prepare the ContainerListOptions
	var containerListOptions types.ContainerListOptions
	if status != "" {
		filterArgs := filters.NewArgs()
		filterArgs.Add("status", status)
		containerListOptions = types.ContainerListOptions{Filters: filterArgs}

	}
	containers, err := Cli.ContainerList(context.Background(), containerListOptions)
	if err != nil {
		return containerList, err
	}

	// Extract containerID, ImageName, and Status
	containerList = make([]Ctn, len(containers))
	for i, container := range containers {
		containerList[i] = Ctn{ID: container.ID[:10], Image: container.Image, Status: container.Status, State: container.State}
	}
	return containerList, nil
}

// ContainersByImageID returns running containers of specific ImageID
func ContainersByImageID(imageID string) (containerList []Ctn, err error) {
	// If containers are filtered by status, prepare the ContainerListOptions
	var containerListOptions types.ContainerListOptions

	filterArgs := filters.NewArgs()
	filterArgs.Add("ancestor", imageID)
	filterArgs.Add("status", "running")

	containerListOptions = types.ContainerListOptions{Filters: filterArgs}
	containers, err := Cli.ContainerList(context.Background(), containerListOptions)
	if err != nil {
		return containerList, err
	}

	// Extract containerID, ImageName, and Status
	containerList = make([]Ctn, len(containers))
	for i, container := range containers {
		containerList[i] = Ctn{ID: container.ID[:10], Image: container.Image, Status: container.Status, State: container.State}
	}
	return containerList, nil
}
