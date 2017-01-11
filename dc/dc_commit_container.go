package dc

import (
	"context"
	"fmt"
	"log"

	"github.com/andreas-kokkalis/dock-server/conf"
	"github.com/docker/docker/api/types"
)

// CommitContainer creates a new image from a running container
func CommitContainer(comment, author, containerID, refTag string) error {

	// TODO: on options, can add a slice of string with the list of changes for this commit
	options := types.ContainerCommitOptions{Comment: comment, Author: author, Reference: conf.GetVal("dc.imagerepo.name") + ":" + refTag}
	response, err := Cli.ContainerCommit(context.Background(), containerID, options)
	if err != nil {
		return err
	}
	// TODO: figure out what to do with the response
	fmt.Printf("%+v\n", response)
	log.Printf("[CommitContainer]: Committed container with ID:%s\n", containerID)
	//sha256:baa8ace946df92b5fb1722538d73531503485535604863e34e174a5d284a601b

	return nil
}
