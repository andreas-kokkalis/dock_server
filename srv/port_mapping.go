package srv

import (
	"errors"
	"fmt"
	"sync"
)

// XXX: This is a very naive implementation.
// One has to consider using channels, and avoid the flags.
// A resource that is free should be part of the data structure,
// A resource that is used, should be removed from the data structure
// TODO: investigate if this is feasible using the chan go idiom

const startPort int = 4200

// PortResources are ports
var PortResources *Resources

// InitPortMappings creates ..
func InitPortMappings(numPorts int) {
	PortResources = &Resources{}
	ReservePorts(PortResources, numPorts)
}

// Resource is a port
type Resource struct {
	port   int
	isUsed bool
}

// Resources are all ports
type Resources struct {
	ports []*Resource
	lock  *sync.Mutex
}

// GetFreeResource is
func GetFreeResource(res *Resources) (int, error) {
	var port int
	res.lock.Lock()

	for _, v := range res.ports {
		if !v.isUsed {
			port = v.port
			v.isUsed = true
			res.lock.Unlock()
			return port, nil
		}
	}
	res.lock.Unlock()
	return port, errors.New("No available port to return")
}

// FreeResource frees a non utilized port
func FreeResource(res *Resources, port int) {
	res.lock.Lock()
	for _, v := range res.ports {
		if v.port == port && v.isUsed == true {
			v.isUsed = false
			res.lock.Unlock()
			return
		}
	}
}

// ReservePorts is not concurrent safe
// TODO: should actually do a list running containers, check if ports are available and then add their values
func ReservePorts(res *Resources, numPorts int) {
	res.lock = &sync.Mutex{}
	myPorts := make([]*Resource, numPorts)

	for i := 0; i < numPorts; i++ {
		myPorts[i] = &Resource{port: startPort + i, isUsed: false}
	}
	res.ports = myPorts
}

// PrintUsed will print the used ports, duh
func PrintUsed(res *Resources) {
	res.lock.Lock()
	for _, v := range res.ports {
		if v.isUsed == true {
			fmt.Printf("Port: %d is used.\n", v.port)
		}
	}
	res.lock.Unlock()
	return
}
