package dc

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/andreas-kokkalis/dock-server/srv"
)

const (
	usrPrefix    = "usr:"
	userTTL      = time.Hour * 1
	admPrefix    = "adm:"
	admRunPrefix = "run:"
	adminTTL     = time.Hour * 24
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

	r, err := GetRunConfig(GetUserKey(userID))
	if err != nil {
		// TODO: parse error
	}
	delPort(r.Port)

	_, err = srv.RCli.Del(GetUserKey(userID)).Result()
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

func setPort(port string, TTL time.Duration) {
	_, err := srv.RCli.Set("port:"+port, true, TTL).Result()
	if err != nil {
		//TODO: do not ignore this error.
	}
}

func delPort(port string) {
	_, err := srv.RCli.Del("port:" + port).Result()
	if err != nil {
		//TODO: do not ignore this error.
	}
}

// ExistsPort is a used by the periodicCheck function to determine whether a running container should be killed, if the corresponding port key has expired.
func ExistsPort(port int) bool {
	exists, _ := srv.RCli.Exists("port:" + strconv.Itoa(port)).Result()
	return exists
}

// GetRunConfig returns the user session
func GetRunConfig(userID string) (r RunConfig, err error) {
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
func SetRunConfig(userID string, r RunConfig) (err error) {
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
	setPort(r.Port, userTTL)
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

/* ===========================================
	Admin Run Container Session
==============================================*/

// GetAdminRunKey constructs the user key
func GetAdminRunKey(key string) string {
	return admRunPrefix + key
}

// DeleteAdminRunConfig deletes the user session
func DeleteAdminRunConfig(key string) error {
	r, err := GetAdminRunConfig(key)
	if err != nil {
		// TODO: parse error
	}
	delPort(r.Port)
	_, err = srv.RCli.Del(GetAdminRunKey(key)).Result()
	if err != nil {
		return err
	}
	return nil
}

// ExistsAdminRunConfig returns true if there is a session for the particular user
func ExistsAdminRunConfig(key string) (bool, error) {
	keyExists, err := srv.RCli.Exists(GetAdminRunKey(key)).Result()
	if err != nil {
		return false, err
	}
	if !keyExists {
		return false, nil
	}
	return true, nil
}

// GetAdminRunConfig returns the user session
func GetAdminRunConfig(key string) (r RunConfig, err error) {
	var val string
	val, err = srv.RCli.Get(GetAdminRunKey(key)).Result()
	if err != nil {
		return r, err
	}
	err = json.Unmarshal([]byte(val), &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// SetAdminRunConfig will add the session
func SetAdminRunConfig(key string, r RunConfig) (err error) {
	// Marshal to JSON
	var js []byte
	js, err = json.Marshal(r)
	if err != nil {
		return err
	}

	// Set key value
	var OK string
	OK, err = srv.RCli.Set(GetAdminRunKey(key), js, userTTL).Result()
	if err != nil {
		return err
	}
	if OK != "OK" {
		return errors.New("Not OK")
	}
	setPort(r.Port, userTTL)
	return nil
}
