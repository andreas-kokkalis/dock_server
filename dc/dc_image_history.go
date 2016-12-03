package dc

import (
	"context"

	"github.com/docker/docker/api/types"
)

// ImageHistory returns image history
func ImageHistory(imageID string) ([]ImgHistory, error) {
	var history []types.ImageHistory
	var err error

	history, err = Cli.ImageHistory(context.Background(), imageID)
	if err != nil {
		return nil, err
	}

	res := []ImgHistory{}

	for _, v := range history {
		res = append(res, ImgHistory{ID: v.ID, Created: v.Created, CreatedBy: v.CreatedBy, Tags: v.Tags, Comment: v.Comment})
		// Only the first is needed
		break
	}

	return res, nil
}
