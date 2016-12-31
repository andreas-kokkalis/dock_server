package dc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// XXX: This is a very naive implementation.
// One has to consider using channels, and avoid the flags.
// A resource that is free should be part of the data structure,
// A resource that is used, should be removed from the data structure
// TODO: investigate if this is feasible using the chan go idiom

const startPort int = 4200

// ContainerPorts are ports
var ContainerPorts *portResources

// ContainerPortsInitialize creates ..
func ContainerPortsInitialize(numPorts int) {
	// Initialize ContainerPorts struct
	ContainerPorts = &portResources{}

	// Lock ContainerPorts and start creating the portResource slice
	ContainerPorts.portsAvailable = make(map[int]bool)
	for i := 0; i < numPorts; i++ {
		ContainerPorts.portsAvailable[startPort+i] = false
	}
}

type portResources struct {
	portsAvailable map[int]bool
	lock           sync.Mutex
}

func (res *portResources) Remove(port int) {
	res.lock.Lock()
	delete(res.portsAvailable, port)
	fmt.Printf("Removed port %d\n", port)
	fmt.Printf("[PortMapper]: removed unused port : %d\n", port)

	res.lock.Unlock()
}

func (res *portResources) fixup(ports map[int]string) {
	res.lock.Lock()
	for port := range res.portsAvailable {
		if _, exists := ports[port]; exists {
			res.portsAvailable[port] = true
		} else {
			res.portsAvailable[port] = false
		}
	}
	res.lock.Unlock()
}

func (res *portResources) Reserve() (port int, err error) {
	res.lock.Lock()
	for port, isUsed := range res.portsAvailable {
		if isUsed == false {
			res.portsAvailable[port] = true
			res.lock.Unlock()
			return port, nil
		}
	}
	res.lock.Unlock()
	return port, errors.New("No available port to return")
}

// PrintUsed will print the used ports, duh
func (res *portResources) PrintUsed() {
	res.lock.Lock()
	for port, isUsed := range res.portsAvailable {
		if isUsed {
			fmt.Printf("[PortMapper]: Port: %d is used.\n", port)
		}
	}
	res.lock.Unlock()
	return
}

// PeriodicChecker checks every X seconds for inconsistencies
// First it gets all used ports by running containers, and syncs the concurrent ports map
// Then it checks if redis configurations exists for the corresponding ports. If such configurations are absent, it will request to kill the containers
func PeriodicChecker() {

	for {
		time.Sleep(time.Second * 3)

		ports, err := GetContainerPorts()
		if err != nil {
			continue
		}

		ContainerPorts.fixup(ports)
		for port, containerID := range ports {
			if !ExistsPort(port) {
				RemoveContainer(containerID, port)
				fmt.Printf("[PortMapper]: removing expired container with ID: %s\n", containerID)
			}
		}
	}
}

// GetContainerPorts returns the list of used ports
func GetContainerPorts() (ports map[int]string, err error) {
	// If containers are filtered by status, prepare the ContainerListOptions
	var containerListOptions types.ContainerListOptions

	filterArgs := filters.NewArgs()
	for _, imageRepo := range GetRepositories() {
		filterArgs.Add("ancestor", imageRepo)
	}
	filterArgs.Add("status", "running")
	containerListOptions = types.ContainerListOptions{Filters: filterArgs}

	containers, err := Cli.ContainerList(context.Background(), containerListOptions)
	if err != nil {
		return ports, err
	}

	// Extract containerID, ImageName, and Status
	ports = make(map[int]string)
	for _, container := range containers {
		fmt.Printf("%+v\n", container.Ports[0].PublicPort)
		// containerList[i] = Ctn{ID: container.ID[:10], Image: container.Image, Status: container.Status, State: container.State}
		ports[int(container.Ports[0].PublicPort)] = container.ID[:10]
	}
	return ports, nil
}
