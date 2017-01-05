package dc

import (
	"context"
	"errors"
	"log"
	"strconv"
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
		ContainerPorts.portsAvailable[startPort+i] = true
	}
}

type portResources struct {
	portsAvailable map[int]bool
	lock           sync.Mutex
}

// Remove is used by dc_remove_container.go and dc_create_container.go to remove unused ports
func (res *portResources) Remove(port int) {
	res.lock.Lock()

	res.portsAvailable[port] = true
	log.Printf("[PortMapper]: removed unused port : %d\n", port)

	res.lock.Unlock()
}

// Reserve is used by dc_create_container.go to reserve an available port
func (res *portResources) Reserve() (port int, err error) {
	res.lock.Lock()

	for port, isAvailable := range res.portsAvailable {
		if isAvailable == true {
			res.portsAvailable[port] = false
			res.lock.Unlock()
			return port, nil
		}
	}
	res.lock.Unlock()
	return port, errors.New("No available port to return")
}

func (res *portResources) fixup(ports map[int]string) {
	res.lock.Lock()

	for port := range res.portsAvailable {
		if _, exists := ports[port]; exists {
			// Reserve port in memory
			res.portsAvailable[port] = false
		} else {
			// Make port available in memory
			res.portsAvailable[port] = true
			// Remove trailing redis configuration
			RemoveIncosistentRedisKeys(strconv.Itoa(port))
		}
	}

	// No ports were used by containers, make sure that none is in memory
	if len(ports) == 0 {
		for port := range res.portsAvailable {
			// Make port available in memory
			res.portsAvailable[port] = true
			// Remove trailing redis configuration
			RemoveIncosistentRedisKeys(strconv.Itoa(port))
		}
	}

	res.lock.Unlock()
}

/*
// PrintUsed will print the used ports, duh
func (res *portResources) PrintUsed() {
	res.lock.Lock()
	for port, isUsed := range res.portsAvailable {
		if isUsed {
			log.Printf("[PortMapper]: Port: %d is used.\n", port)
		}
	}
	res.lock.Unlock()
	return
}
*/

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

		// Check for containers that have crashed / stopped etc.
		// Remove the PortsAvailable
		// Remove their redis keys
		ContainerPorts.fixup(ports)

		// Check for expired redis keys
		for port, containerID := range ports {
			if !ExistsPort(strconv.Itoa(port)) {
				RemoveContainer(containerID, port)
				log.Printf("[PortMapper]: removing expired container with ID: %s and port: %d\n", containerID, port)
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
		log.Printf("[PortMapper]: port %v is in use by contaienr %v\n", container.Ports[0].PublicPort, container.ID[:12])
		// containerList[i] = Ctn{ID: container.ID[:10], Image: container.Image, Status: container.Status, State: container.State}
		ports[int(container.Ports[0].PublicPort)] = container.ID[:12]
	}
	return ports, nil
}
