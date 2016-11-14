package dc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// RemoveImage removes an image
func RemoveImage(imageID string) error {

	options := types.ImageRemoveOptions{}
	imgDelete, err := Cli.ImageRemove(context.Background(), imageID, options)
	if err != nil {
		return err
	}
	// TODO: figure out what to do with the imgDelete
	fmt.Println(imgDelete)

	return nil
}
