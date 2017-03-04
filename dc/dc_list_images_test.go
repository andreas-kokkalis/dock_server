package dc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListImages(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	InitTestDependencies()
	imgList, err := ListImages()
	assert.NotNil(imgList, "It should not be nil")
	assert.NoError(err, "It should not return an error")

	for _, img := range imgList {
		assert.True(vTestImageID.MatchString(img.ID), "It should be true")
		assert.NotEqual(0, len(img.RepoTags), "It should not have 0 length")
	}
}

func TestGetRepositories(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	repos := GetRepositories()
	assert.NotEqual(0, len(repos))
	for _, repo := range repos {
		fmt.Println(repo)
	}
}
