package dc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// CommitContainer creates a new image from a running container
func CommitContainer(comment, author, containerID, refTag string) error {

	// TODO: on options, can add a slice of string with the list of changes for this commit
	options := types.ContainerCommitOptions{Comment: comment, Author: author, Reference: refTag}
	response, err := Cli.ContainerCommit(context.Background(), containerID, options)
	if err != nil {
		return err
	}
	// TODO: figure out what to do with the response
	fmt.Println(response)
	//sha256:baa8ace946df92b5fb1722538d73531503485535604863e34e174a5d284a601b

	return nil
}
