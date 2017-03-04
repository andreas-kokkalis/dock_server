package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Parallel()
	name := "correct path"
	err := InitConf("./")
	assert.NoError(t, err, name)
}

type confVals struct {
	key, val string
}

var vals = []confVals{
	{"dc.imagerepo.name", "dc"},
	{"dc.docker.api.host", "unix:///var/run/docker.sock"},
	{"dc.docker.api.version", "1.24"},
}

func TestGetVal(t *testing.T) {
	t.Parallel()
	InitConf("./")
	for _, test := range vals {
		val := GetVal(test.key)
		if val != test.val {
			t.Errorf("Conf key: %s, Expected: %s Got: %s", test.key, test.val, val)
		}
	}
}
