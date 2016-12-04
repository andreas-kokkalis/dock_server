package session

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/srv"
)

const (
	usrPrefix = "usr:"
	userTTL   = time.Hour * 1
	admPrefix = "adm:"
	adminTTL  = time.Hour * 24
)

/*
	============================
				USER
	============================
*/

// GetUserKey constructs the user key
func GetUserKey(userID string) string {
	return usrPrefix + userID
}

func StripUserKey(key string) string {
	return key[4:]
}

// DeleteRunConfig deletes the user session
func DeleteRunConfig(userID string) error {
	_, err := srv.RCli.Del(GetUserKey(userID)).Result()
	if err != nil {
		return err
	}
	return nil
}

// ExistsRunConfig returns true if there is a session for the particular user
func ExistsRunConfig(userID string) (bool, error) {
	keyExists, err := srv.RCli.Exists(GetUserKey(userID)).Result()
	if err != nil {
		return false, err
	}
	if !keyExists {
		return false, nil
	}
	return true, nil
}

// GetRunConfig returns the user session
func GetRunConfig(userID string) (r dc.RunConfig, err error) {
	var val string
	val, err = srv.RCli.Get(GetUserKey(userID)).Result()
	if err != nil {
		return r, err
	}
	err = json.Unmarshal([]byte(val), &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SetRunConfig will add the session
func SetRunConfig(userID string, r dc.RunConfig) (err error) {
	// Marshal to JSON
	var js []byte
	js, err = json.Marshal(r)
	if err != nil {
		return err
	}

	// Set key value
	var OK string
	OK, err = srv.RCli.Set(GetUserKey(userID), js, userTTL).Result()
	if err != nil {
		return err
	}
	if OK != "OK" {
		return errors.New("Not OK")
	}
	return nil
}

/*
	============================
				ADMIN
	============================
*/

type adminSession struct {
	AdminID int `json:"id"`
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
