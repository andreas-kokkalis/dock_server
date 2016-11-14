package dc

import "context"

// KillContainer kills a container
func KillContainer(containerID string) error {
	err := Cli.ContainerKill(context.Background(), containerID, "SIGKILL")
	if err != nil {
		return err
	}
	return nil
}
