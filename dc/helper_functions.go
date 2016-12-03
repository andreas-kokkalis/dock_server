package dc

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
)

// GetTagByID returns an imageTag given and imageID:
// 	performs an ImageList request, gathers all results and returns the tag of the given imageID
func GetTagByID(imageID string) (string, error) {
	images, err := Cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return "", err
	}
	// Extract imageID, RepoTags for specific type of images
	for _, image := range images {
		if image.ID[7:19] == imageID {
			s := image.RepoTags[0]
			return s[strings.LastIndex(s, ":")+1:], nil
		}
	}
	return "", nil
}
