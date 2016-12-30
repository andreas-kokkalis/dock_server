package dc

import (
	"context"
	"fmt"
	"strings"

	"github.com/andreas-kokkalis/dock-server/conf"

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
	// XXX: consider limiting this to only the original base image. ==> NOT SUPPORTED

	// the args will be {"image.name":{"ubuntu":true},"label":{"label1=1":true,"label2=2":true}}
	//map[string]map[string]bool

	// From the docker daemon source code:
	/*var acceptedImageFilterTags = map[string]bool{
		"dangling":  true,
		"label":     true,
		"before":    true,
		"since":     true,
		"reference": true,
	}
	*/

	images, err := Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return imageList, err
	}

	// Extract imageID, RepoTags for specific type of images
	for _, image := range images {
		fmt.Println(image.RepoTags[0])
		s := image.RepoTags[0]
		if s[0:strings.LastIndex(s, ":")] == conf.GetVal("dc.imagerepo.name") {
			imageList = append(imageList, Img{ID: image.ID[7:19], RepoTags: image.RepoTags})
		}
	}
	return imageList, nil
}

// GetRepositories returns the lsit of the dc repositories and tags
func GetRepositories() (imageList string) {
	images, err := Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return imageList
	}

	// Extract imageID, RepoTags for specific type of images

	for _, image := range images {
		fmt.Println(image.RepoTags[0])
		s := image.RepoTags[0]
		if s[0:strings.LastIndex(s, ":")] == conf.GetVal("dc.imagerepo.name") {
			imageList = imageList + " " + image.RepoTags[0]
		}
	}
	return imageList
}
