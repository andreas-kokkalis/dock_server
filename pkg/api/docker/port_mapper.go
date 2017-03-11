package docker

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/andreas-kokkalis/dock-server/pkg/api/store"
)

const startPort int = 4200

// XXX: This is a very naive implementation.
// One has to consider using channels, and avoid the flags.
// A resource that is free should be part of the data structure,
// A resource that is used, should be removed from the data structure
// TODO: investigate if this is feasible using the chan go idiom

// TODO: Possibly make portResources part of PortMapper

// PortMapper ...
type PortMapper struct {
	redis *store.RedisRepo
	ports *portResources
}

type portResources struct {
	portsAvailable map[int]bool
	lock           sync.Mutex
}

// NewPortMapper initializes the port mapper
func NewPortMapper(redis *store.RedisRepo, numPorts int) *PortMapper {
	mapper := &PortMapper{redis: redis}
	mapper.ports = &portResources{}
	mapper.ports.portsAvailable = make(map[int]bool)
	for i := 0; i < numPorts; i++ {
		mapper.ports.portsAvailable[startPort+i] = true
	}
	return mapper
}

// Remove is used by dc_remove_container.go and dc_create_container.go to remove unused ports
func (pm *PortMapper) Remove(port int) {
	pm.ports.lock.Lock()

	pm.ports.portsAvailable[port] = true
	log.Printf("[PortMapper]: removed unused port : %d\n", port)

	pm.ports.lock.Unlock()
}

// Reserve is used by dc_create_container.go to reserve an available port
func (pm *PortMapper) Reserve() (port int, err error) {
	pm.ports.lock.Lock()

	for port, isAvailable := range pm.ports.portsAvailable {
		if isAvailable == true {
			pm.ports.portsAvailable[port] = false
			pm.ports.lock.Unlock()
			return port, nil
		}
	}
	pm.ports.lock.Unlock()
	return port, errors.New("No available port to return")
}

func (pm *PortMapper) fixup(ports map[int]string) {
	pm.ports.lock.Lock()

	for port := range pm.ports.portsAvailable {
		if _, exists := ports[port]; exists {
			// Reserve port in memory
			pm.ports.portsAvailable[port] = false
		} else {
			// Make port available in memory
			pm.ports.portsAvailable[port] = true
			// Remove trailing redis configuration
			pm.redis.RemoveIncosistentRedisKeys(strconv.Itoa(port))
		}
	}

	// No ports were used by containers, make sure that none is in memory
	if len(ports) == 0 {
		for port := range pm.ports.portsAvailable {
			// Make port available in memory
			pm.ports.portsAvailable[port] = true
			// Remove trailing redis configuration
			pm.redis.RemoveIncosistentRedisKeys(strconv.Itoa(port))
		}
	}

	pm.ports.lock.Unlock()
}

// // PrintUsed will print the used ports, duh
// func (res *portResources) PrintUsed() {
// 	res.lock.Lock()
// 	for port, isUsed := range res.portsAvailable {
// 		if isUsed {
// 			log.Printf("[PortMapper]: Port: %d is used.\n", port)
// 		}
// 	}
// 	res.lock.Unlock()
// 	return
// }

// PeriodicChecker checks every X seconds for inconsistencies
// First it gets all used ports by running containers, and syncs the concurrent ports map
// Then it checks if redis configurations exists for the corresponding ports. If such configurations are absent, it will request to kill the containers
func PeriodicChecker(docker *Repo, pm *PortMapper, redis *store.RedisRepo) {

	for {
		time.Sleep(time.Second * 3)

		ports, err := docker.ContainerGetUsedPorts()
		if err != nil {
			continue
		}

		// Check for containers that have crashed / stopped etc.
		// Remove the PortsAvailable
		// Remove their redis keys
		pm.fixup(ports)

		// Check for expired redis keys
		for port, containerID := range ports {
			if !redis.ExistsPort(strconv.Itoa(port)) {
				_ = docker.ContainerRemove(containerID, port)
				pm.Remove(port)

				log.Printf("[PortMapper]: removing expired container with ID: %s and port: %d\n", containerID, port)
			}
		}
	}
}
