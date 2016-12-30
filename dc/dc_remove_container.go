package dc

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
)

// RemoveContainer force removes a container
func RemoveContainer(containerID string, port int) (err error) {
	t := time.Duration(time.Millisecond * 100)

	// First stop the container
	err = Cli.ContainerStop(context.Background(), containerID, &t)
	if err != nil {
		// shut up ...
		fmt.Println("Attemted to stop the container")
		//return err
	}

	// Then kill it
	err = Cli.ContainerKill(context.Background(), containerID, "SIGKILL")
	if err != nil {
		fmt.Println("Attemted to kill the container")
		// srv.FreeResource(srv.PortResources, port)
		// return err
	}

	// After the container is killed free the port resource
	ContainerPorts.Remove(port)
	// rm -f the container
	options := types.ContainerRemoveOptions{Force: true}
	err = Cli.ContainerRemove(context.Background(), containerID, options)
	if err != nil {
		fmt.Println("attempted to remove the container")
		fmt.Println(err.Error())
		// return err
	}
	return nil
}
