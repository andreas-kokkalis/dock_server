package docker

import (
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TODO: mock client connection and perform tests that do not rely on a docker daemon running.
// Move these tests to the integration testing

func getRepo() *Repo {
	c, _ := config.NewConfig(validConfigDir, "local")
	d, _ := NewAPIClient(c.GetDockerConfig())
	return NewRepo(d, c.GetDockerConfig())
}

func TestNewRepo(t *testing.T) {
	t.Parallel()
	r := getRepo()
	assert.NotNil(t, r)
}

/*
func TestImageList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	r := getRepo()
	imgList, err := r.ImageList()
	assert.NoError(err)
	assert.NotEqual(0, len(imgList))
}

func TestImageListRepositories(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	r := getRepo()
	repoList := r.ImageListRepositories()
	assert.NotEqual(0, len(repoList))
}

func TestImageTagByID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	r := getRepo()
	tag, err := r.ImageTagByID("")
	assert.NoError(err)
	assert.Empty(tag)

	imgList, _ := r.ImageList()
	tag, err = r.ImageTagByID(imgList[0].ID)
	assert.NoError(err)
	assert.NotEmpty(tag)
}

func TestContainerCheckState(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	r := getRepo()
	// Returns error cause containerID is invalid
	isRunning, err := r.ContainerCheckState("", "running")
	assert.Error(err)
	assert.Equal(false, isRunning)

	// TODO:
	// Run a container
	// Check state
}
*/
