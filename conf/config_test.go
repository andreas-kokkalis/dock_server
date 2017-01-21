package conf

import "testing"

func TestInit(t *testing.T) {
	err := InitConf("./conf")
	if err == nil {
		t.Errorf("Expected Error, got nil")
	}
	err = InitConf("./")
	if err != nil {
		t.Errorf("Expected nil error, got %s", err.Error())
	}
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
	InitConf("./")
	for _, test := range vals {
		val := GetVal(test.key)
		if val != test.val {
			t.Errorf("Conf key: %s, Expected: %s Got: %s", test.key, test.val, val)
		}
	}
}
