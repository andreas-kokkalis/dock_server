package dc

import "time"

// Ctn minimal container struct
type Ctn struct {
	ID     string `json:"Id"`
	Image  string
	Status string
	State  string
}

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

// RunConfig for running a container
type RunConfig struct {
	ContainerID string `json:"id"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	URL         string `json:"url"`
}
