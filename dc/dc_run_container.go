package dc

import (
	"log"
	"strconv"
)

// RunContainer does something
func RunContainer(imageID, username, password string) (cfg RunConfig, err error) {
	// Create the container
	id, port, err := CreateContainer(imageID, password)
	if err != nil {
		log.Printf("[RunContainer]: Error while creating: %v\n", err.Error())
		return cfg, err
	}
	// Start the container
	err = StartContainer(id)
	if err != nil {
		log.Printf("[RunContainer]: Error while starting: %v\n", err.Error())
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
