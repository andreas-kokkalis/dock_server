package api

import (
	"regexp"
	"time"
)

// VImageID is a validator for image Identifier
var VImageID = regexp.MustCompile(`^([A-Fa-f0-9]{12,64})$`)

// Img minimal image struct
type Img struct {
	ID        string    `json:"Id"`
	RepoTags  []string  `json:"RepoTags"`
	CreatedAt time.Time `json:"CreatedAt"`
}

// ImgHistory is identical to docker types.ImgHistory
type ImgHistory struct {
	ID        string    `json:"Id"`
	CreatedAt time.Time `json:"CreatedAt"`
	CreatedBy string    `json:"CreatedBy"`
	RepoTags  []string  `json:"RepoTags"`
	Size      int64     `json:"Size"`
	Comment   string    `json:"Comment"`
}
