package api

// Ctn minimal container struct
type Ctn struct {
	ID     string `json:"Id"`
	Image  string
	Status string
	State  string
}
