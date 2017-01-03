package dc

import (
	"fmt"
	"strconv"
)

// RunContainer does something
func RunContainer(imageID, username, password string) (cfg RunConfig, err error) {
	// Create the container
	id, port, err := CreateContainer(imageID, username, password)
	if err != nil {
		fmt.Printf("error-create: %v\n", err.Error())
		return cfg, err
	}
	// Start the container
	err = StartContainer(id)
	if err != nil {
		fmt.Printf("error-start: %v\n", err.Error())
		return cfg, err
	}

	cfg = RunConfig{
		ContainerID: id,
		Username:    username,
		Password:    password,
		Port:        strconv.Itoa(port),
		URL:         "https://127.0.0.1:" + strconv.Itoa(port),
	}
	return cfg, nil
}
