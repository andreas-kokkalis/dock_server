package session

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/andreas-kokkalis/dock-server/srv"
)

const (
	usrPrefix = "dsx:"
	admPrefix = "adm:"
	adminTTL  = time.Hour * 24
)

type userSession struct {
	ContainerID string `json:"ID"`
	Port        int    `json:"port"`
}

type adminSession struct {
	AdminID int `json:"id"`
}

// UserExists returns true if there is a session for the particular user
func UserExists(username, password string) (bool, error) {
	key := usrPrefix + username + ":" + password

	keyExists, err := srv.RCli.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if !keyExists {
		return false, nil
	}
	return true, nil
}

// UserAdd will add the session
func UserAdd(username, password, containerID string, port int, ttl int) error {
	key := usrPrefix + username + ":" + password

	exists, err := UserExists(username, password)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Key already exists")
	}

	s := userSession{ContainerID: containerID, Port: port}
	OK, err2 := srv.RCli.Set(key, s, time.Duration(ttl)).Result()
	if err2 != nil {
		return err2
	}
	if OK != "OK" {
		return errors.New(OK)
	}
	fmt.Println(OK)
	return nil
}

// GetAdminKey returns the admin session key
func GetAdminKey(adminID int) string {
	h := md5.New()
	io.WriteString(h, strconv.Itoa(adminID))
	io.WriteString(h, "key")
	s := fmt.Sprintf("%x", h.Sum(nil))
	return admPrefix + s

}

// AdminExists checks if a session exists for that particular adminID
func AdminExists(adminID int) (bool, error) {
	exists, err := srv.RCli.Exists(GetAdminKey(adminID)).Result()
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, nil
}

// AdminAdd will add a key for that admin
func AdminAdd(adminID int) error {
	key := GetAdminKey(adminID)
	err := srv.RCli.Set(key, adminID, 0).Err()

	time.Sleep(1000)
	_, err = srv.RCli.Exists(GetAdminKey(adminID)).Result()
	if err != nil {
		return err
	}
	return nil
}
