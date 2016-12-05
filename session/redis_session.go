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

// StripKey removes the prefix
func StripKey(key string) string {
	return key[4:]
}

/* ============================
				USER
   ============================ */

// GetUserKey constructs the user key
func GetUserKey(userID string) string {
	return usrPrefix + userID
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

// CreateAdminKey returns the admin session key
func CreateAdminKey(adminID int) string {
	h := md5.New()
	io.WriteString(h, strconv.Itoa(adminID))
	io.WriteString(h, "key")
	s := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(admPrefix + s)
	return admPrefix + s
}

// ExistsAdminSession checks if a session exists for that particular adminID
func ExistsAdminSession(key string) (bool, error) {
	exists, err := srv.RCli.Exists(key).Result()
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, nil
}

// SetAdminSession will add a key for that admin
func SetAdminSession(key string) error {
	err := srv.RCli.Set(key, key, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// DeleteAdminSession will add a key for that admin
func DeleteAdminSession(key string) error {
	_, err := srv.RCli.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}
