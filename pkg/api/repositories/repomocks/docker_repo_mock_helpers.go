package repomocks

// NewDockerRepositoryMock initializes a mock implementations of DockerRepository
func NewDockerRepositoryMock() *DockerRepositoryMock {
	return &DockerRepositoryMock{}
}

// WithContainerGetUsedPorts sets a mock function of ContainerGetUsedPorts
func (d *DockerRepositoryMock) WithContainerGetUsedPorts(ports map[int]string, err error) *DockerRepositoryMock {
	d.ContainerGetUsedPortsFunc = func() (map[int]string, error) {
		return ports, err
	}
	return d
}

// WithContainerRemove sets the ContainerRemove mock function
func (d *DockerRepositoryMock) WithContainerRemove(err error) *DockerRepositoryMock {
	d.ContainerRemoveFunc = func(_ string, _ int) error {
		return err
	}
	return d
}
