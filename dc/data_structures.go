package dc

// Ctn minimal container struct
type Ctn struct {
	ID     string `json:"Id"`
	Image  string
	Status string
	State  string
}

// Img minimal image struct
type Img struct {
	ID       string   `json:"Id"`
	RepoTags []string `json:"RepoTags"`
}

// ImgHistory is identical to docker types.ImgHistory
type ImgHistory struct {
	ID        string `json:"Id"`
	Created   int64
	CreatedBy string
	Tags      []string
	Size      int64
	Comment   string
}
