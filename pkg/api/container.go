package api

import "regexp"

// VContainerID calidates if a container id is valid
var VContainerID = regexp.MustCompile(`^([A-Fa-f0-9]{12,64})$`)

// VContainerState ...
var VContainerState = regexp.MustCompile(`^(|created|restarting|running|paused|exited|dead)$`) // Can also be empty

// Ctn minimal container struct
type Ctn struct {
	ID     string `json:"Id"`
	Image  string
	Status string
	State  string
}

// ContainerRun models responses for a running container
type ContainerRun struct {
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ContainerID string `json:"id"`
}

// ContainerCommitResponse models the response of creating a new image from a running container.
type ContainerCommitResponse struct {
	ImageID string `json:"imageID"`
}

// ContainerCommitRequest models a request for creating an image from a runnign container
type ContainerCommitRequest struct {
	Comment string `json:"comment"`
	Author  string `json:"auth"`
	RefTag  string `json:"tag"`
}
