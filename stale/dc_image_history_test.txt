package dc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageHistory(t *testing.T) {
	assert := assert.New(t)

	InitTestDependencies()
	imgList, _ := ListImages()

	history, err2 := ImageHistory(imgList[0].ID)
	assert.NotNil(history, "It should not be nil")
	assert.NoError(err2, "It should not return error")

	assert.True(vTestImageID.MatchString(history[0].ID), "It should be true")
	assert.NotEqual(t, 0, len(history[0].RepoTags), "It should not have 0 length")
}
