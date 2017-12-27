package repositories

import (
	"github.com/andreas-kokkalis/dock_server/pkg/api"
)

//go:generate moq -out ../repomocks/docker_repo_mock.go -pkg repomocks . DockerRepository

// DockerRepository models the interaction with Docker daemon for the purposes of the dock_server API
type DockerRepository interface {
	ContainerGetUsedPorts() (map[int]string, error)
	ContainerRemove(containerID string, port int) error
	ContainerRun(imageID, username, password string, port int) (api.RunConfig, error)
	ContainerCommit(comment, author, containerID, refTag string) (string, error)
	ContainerList(status string) ([]api.Ctn, error)
	ImageList() ([]api.Img, error)
	ImageHistory(imageID string) ([]api.ImgHistory, error)
	ImageRemove(imageID string) error
	ImageGetTagByID(imageID string) (string, error)
	GetRunningContainersByImageID(imageID string) ([]api.Ctn, error)
}
