package portmapper

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories/repomocks"
)

func TestNewPortMapper(t *testing.T) {
	m := repomocks.NewRedisRepositoryMock()
	pm := NewPortMapper(m, 10)
	assert.NotNil(t, pm.ports)
	assert.Equal(t, 10, len(pm.ports.portsAvailable))
}

func TestReserveRemove(t *testing.T) {
	pm := NewPortMapper(nil, 1)
	port, err := pm.Reserve()
	assert.NoError(t, err, "reserve port 0/1")
	assert.Equal(t, 4200, port)

	port2, err := pm.Reserve()
	assert.Error(t, err, "reserve port 1/1")
	assert.Equal(t, 0, port2)

	pm.Remove(port)
}

func TestFixup(t *testing.T) {

	redis := repomocks.NewRedisRepositoryMock().
		WithDeleteStaleMappedPort()

	pm := NewPortMapper(redis, 5)
	ports := make(map[int]string)

	port1, err := pm.Reserve()
	assert.NoError(t, err, "reserve port 0/1")
	assert.NotEqual(t, 0, port1)
	ports[port1] = "port1"

	port2, err := pm.Reserve()
	assert.NoError(t, err, "reserve port 0/1")
	assert.NotEqual(t, 0, port2)
	ports[port2] = "port2"

	// Missing port in the mapper should be adde by fixup
	port3, _ := pm.Reserve()
	pm.Remove(port3)
	ports[port3] = "port3"
	assert.True(t, pm.ports.portsAvailable[port3], "port3 is not yet reserved")
	pm.fixup(ports)
	assert.False(t, pm.ports.portsAvailable[port3], "port3 is now reserved")

	// fixup empty map will clear all ports
	pm.fixup(make(map[int]string))
	for _, available := range pm.ports.portsAvailable {
		assert.True(t, available)
	}
}

func TestCheck(t *testing.T) {

	// First scenario, ports are reserved, Getting ports of containers returns error
	// Do not remove any used ports
	emptyPortsMap := make(map[int]string)
	redis := repomocks.NewRedisRepositoryMock().
		WithPortIsMapped(true).
		WithDeleteStaleMappedPort()
	docker := repomocks.NewDockerRepositoryMock().
		WithContainerRemove(nil).
		WithContainerGetUsedPorts(emptyPortsMap, errors.New("No results"))

	pm := NewPortMapper(redis, 5)
	port1, err := pm.Reserve()
	assert.NoError(t, err, "reserve port 0/1")
	assert.NotEqual(t, 0, port1)

	port2, err := pm.Reserve()
	assert.NoError(t, err, "reserve port 0/1")
	assert.NotEqual(t, 0, port2)

	Check(docker, pm, redis)
	for port, available := range pm.ports.portsAvailable {
		if port == port1 || port == port2 {
			assert.False(t, available, fmt.Sprintf("port: %d", port))
		} else {
			assert.True(t, available, fmt.Sprintf("port: %d", port))
		}
	}

	// second scenario 2 containers are running, but nothing is reserved in redis. Delete containers and removed reservations.
	portsMap := make(map[int]string)
	pm = NewPortMapper(redis, 5)
	port1, err = pm.Reserve()
	assert.NoError(t, err, "reserve port 0/1")
	assert.NotEqual(t, 0, port1)
	portsMap[port1] = "port1"

	port2, err = pm.Reserve()
	assert.NoError(t, err, "reserve port 0/1")
	assert.NotEqual(t, 0, port2)
	portsMap[port2] = "port2"

	redis = repomocks.NewRedisRepositoryMock().
		WithPortIsMapped(false).
		WithDeleteStaleMappedPort()
	docker = repomocks.NewDockerRepositoryMock().
		WithContainerRemove(nil).
		WithContainerGetUsedPorts(portsMap, nil)

	Check(docker, pm, redis)
	for port, available := range pm.ports.portsAvailable {
		assert.True(t, available, fmt.Sprintf("port: %d", port))
	}
}
