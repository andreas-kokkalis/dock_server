package api

// RunConfig for running a container
type RunConfig struct {
	ContainerID string `json:"id"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	URL         string `json:"url"`
}
