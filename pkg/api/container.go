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
