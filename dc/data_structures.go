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
