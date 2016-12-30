package dc

import (
	"fmt"
	"time"
)

// ResourceHandler verifies if the system resources are behaving as expected
func ResourceHandler() {
	for {
		time.Sleep(time.Second * 2)
		fmt.Println(GetRepositories())

		// Get all resources from Redis

		// Get all resources from Docker daemon
		//      identify all images
		//      get all containers of all images
		//      get all ports of all containers

		// Verify all ports from Concurrent storage
	}
}
