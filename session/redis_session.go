package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/andreas-kokkalis/dock-server/srv"
)

const prefix string = "dsx:"

type sess struct {
	ContainerID string `json:"ID"`
	Port        int    `json:"port"`
}

// Exists returns where there is a session for the particular user
func Exists(username, password string) (bool, error) {
	key := prefix + username + ":" + password

	keyExists, err := srv.RCli.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if !keyExists {
		return false, nil
	}
	return true, nil
}

// Add returns where there is a session for the particular user
func Add(username, password, containerID string, port int, ttl int) error {
	key := prefix + username + ":" + password

	exists, err := Exists(username, password)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Key already exists")
	}

	s := sess{ContainerID: containerID, Port: port}
	keyExists, err2 := srv.RCli.Set(key, s, time.Duration(ttl)).Result()
	if err2 != nil {
		return err2
	}
	fmt.Println(keyExists)
	// if !keyExists {
	// return false, nil
	// }
	return nil
}
