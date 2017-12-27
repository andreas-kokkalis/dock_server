package repomocks

import "github.com/andreas-kokkalis/dock_server/pkg/api"

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

// WithImageList sets the ImageList mock function
func (d *DockerRepositoryMock) WithImageList(images []api.Img, err error) *DockerRepositoryMock {
	d.ImageListFunc = func() ([]api.Img, error) {
		return images, err
	}
	return d
}

// WithImageHistory sets the ImageHistory mock function
func (d *DockerRepositoryMock) WithImageHistory(history []api.ImgHistory, err error) *DockerRepositoryMock {
	d.ImageHistoryFunc = func(_ string) ([]api.ImgHistory, error) {
		return history, err
	}
	return d
}

// WithImageRemove sets the ImageRemove mock function
func (d *DockerRepositoryMock) WithImageRemove(err error) *DockerRepositoryMock {
	d.ImageRemoveFunc = func(_ string) error {
		return err
	}
	return d
}

// WithGetRunningContainersByImageID sets the GetRunningContainersByImageID mock function
func (d *DockerRepositoryMock) WithGetRunningContainersByImageID(containers []api.Ctn, err error) *DockerRepositoryMock {
	d.GetRunningContainersByImageIDFunc = func(_ string) ([]api.Ctn, error) {
		return containers, err
	}
	return d
}
