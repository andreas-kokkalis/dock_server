package dc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// RemoveImage removes an image
func RemoveImage(imageID string) error {

	options := types.ImageRemoveOptions{}
	tag, err := GetTagByID(imageID)
	imgDelete, err := Cli.ImageRemove(context.Background(), imageRepo+":"+tag, options)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// TODO: figure out what to do with the imgDelete
	fmt.Println(imgDelete)

	return nil
}
