package dc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// ListImages retrieves the list of docker images from Docker Engine
// via the Docker Remote API.
// It only returns the image ID, the repotags
func ListImages() ([]Img, error) {
	var imageList []Img

	// Get list of images from Docker Engine
	// types.ImageListOptions accepts filters.
	// Since no filters are used, all images are returned.
	// XXX: consider limiting this to only the original base image.
	images, err := Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return imageList, err
	}

	// Extract imageID, RepoTags
	imageList = make([]Img, len(images))
	for i, image := range images {

		fmt.Println(image)
		imageList[i] = Img{ID: image.ID[7:19], RepoTags: image.RepoTags}
	}
	return imageList, nil
}
