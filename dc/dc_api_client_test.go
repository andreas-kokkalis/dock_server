package dc

import (
	"log"
	"regexp"
	"testing"

	"github.com/andreas-kokkalis/dock-server/conf"
	"github.com/stretchr/testify/assert"
)

func initTestDependencies() {
	// Load static configuration strings from conf/conf.yaml
	err := conf.InitConf("../conf")
	if err != nil {
		log.Fatal(err)
	}
	ContainerPortsInitialize(200)
}

//InitTestDependencies is used for the tests of this package
func InitTestDependencies() {
	// Load static configuration strings from conf/conf.yaml
	err := conf.InitConf("../conf")
	if err != nil {
		log.Fatal(err)
	}
	ContainerPortsInitialize(200)
	APIClientInit(conf.GetVal("dc.docker.api.version"), conf.GetVal("dc.docker.api.host"))
}

func TestAPIClientInit(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		APIClientInit(conf.GetVal("dc.docker.api.version"), conf.GetVal("dc.docker.api.host"))
	}, "It should panic.")

	initTestDependencies()
	assert.NotPanics(func() {
		APIClientInit(conf.GetVal("dc.docker.api.version"), conf.GetVal("dc.docker.api.host"))
	}, "It should panic.")

	assert.NotNil(t, Cli, "It should not be nil")
}

var vTestImageID = regexp.MustCompile(`^([A-Fa-f0-9]{12,64})$`)
