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

	for i, v := range history {
		res[i].ID = v.ID
		res[i].Created = v.Created
		res[i].CreatedBy = v.CreatedBy
		res[i].Tags = v.Tags
	}
	return res, nil
}
