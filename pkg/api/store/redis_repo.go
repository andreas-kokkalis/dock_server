package store

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
)

// RedisRepo ...
type RedisRepo struct {
	redis RedisClient
}

// NewRedisRepo ...
func NewRedisRepo(redis RedisClient) *RedisRepo {
	return &RedisRepo{redis}
}

const (
	usrPrefix    = "usr:"
	userTTL      = time.Hour * 5
	admPrefix    = "adm:"
	admRunPrefix = "run:"
	adminTTL     = time.Hour * 24
)

// StripSessionKeyPrefix removes the prefix
func (r *RedisRepo) StripSessionKeyPrefix(key string) string {
	return key[4:]
}

/* ============================
				USER
   ============================ */

// GetUserRunKey constructs the user key
func (r *RedisRepo) GetUserRunKey(userID string) string {
	return usrPrefix + userID
}

// DeleteUserRunConfig deletes the user session
func (r *RedisRepo) DeleteUserRunConfig(userID string) error {

	runConfig, err := r.GetUserRunConfig(r.GetUserRunKey(userID))
	if err != nil {
		// TODO: parse error
	}
	r.delPort(runConfig.Port)

	_, err = r.redis.Del(r.GetUserRunKey(userID))
	return err
}

// ExistsUserRunConfig returns true if there is a session for the particular user
func (r *RedisRepo) ExistsUserRunConfig(userID string) (bool, error) {
	keyExists, err := r.redis.Exists(r.GetUserRunKey(userID))
	if err != nil {
		return false, err
	}
	if !keyExists {
		return false, nil
	}
	return true, nil
}

// GetUserRunConfig returns the user session
func (r *RedisRepo) GetUserRunConfig(userID string) (runConfig api.RunConfig, err error) {
	var val string
	val, err = r.redis.Get(r.GetUserRunKey(userID))
	if err != nil {
		return runConfig, err
	}
	err = json.Unmarshal([]byte(val), &runConfig)
	return runConfig, err
}

// SetUserRunConfig will add the session
func (r *RedisRepo) SetUserRunConfig(userID string, runConfig api.RunConfig) (err error) {
	// Marshal to JSON
	var js []byte
	js, err = json.Marshal(runConfig)
	if err != nil {
		return err
	}

	// Set key value
	var OK string
	OK, err = r.redis.Set(r.GetUserRunKey(userID), string(js), userTTL)
	if err != nil {
		return err
	}
	if OK != "OK" {
		return errors.New("Not OK")
	}
	r.setPort(runConfig.Port, r.GetUserRunKey(userID), userTTL)
	return nil
}

/*
	============================
				ADMIN
	============================
*/

// CreateAdminKey returns the admin session key
func (r *RedisRepo) CreateAdminKey(adminID int) string {
	h := md5.New()
	_, _ = io.WriteString(h, strconv.Itoa(adminID))
	_, _ = io.WriteString(h, "key")
	s := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(admPrefix + s)
	return admPrefix + s
}

// ExistsAdminSession checks if a session exists for that particular adminID
func (r *RedisRepo) ExistsAdminSession(key string) (bool, error) {
	exists, err := r.redis.Exists(key)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	return true, nil
}

// SetAdminSession will add a key for that admin
func (r *RedisRepo) SetAdminSession(key string) error {
	_, err := r.redis.Set(key, key, 0)
	return err
}

// DeleteAdminSession will add a key for that admin
func (r *RedisRepo) DeleteAdminSession(key string) error {
	_, err := r.redis.Del(key)
	return err
}

/* ===========================================
	Admin Run Container Session
==============================================*/

// GetAdminSessionRunKey constructs the admin run key
func (r *RedisRepo) GetAdminSessionRunKey(key string) string {
	return admRunPrefix + key
}

// DeleteAdminRunConfig deletes the user session
func (r *RedisRepo) DeleteAdminRunConfig(key string) error {
	runConfig, err := r.GetAdminRunConfig(key)
	if err != nil {
		// TODO: parse error
	}
	r.delPort(runConfig.Port)
	_, err = r.redis.Del(r.GetAdminSessionRunKey(key))
	return err
}

// ExistsAdminRunConfig returns true if there is a session for the particular user
func (r *RedisRepo) ExistsAdminRunConfig(key string) (bool, error) {
	keyExists, err := r.redis.Exists(r.GetAdminSessionRunKey(key))
	if err != nil {
		return false, err
	}
	if !keyExists {
		return false, nil
	}
	return true, nil
}

// GetAdminRunConfig returns the user session
func (r *RedisRepo) GetAdminRunConfig(key string) (runConfig api.RunConfig, err error) {
	var val string
	val, err = r.redis.Get(r.GetAdminSessionRunKey(key))
	if err != nil {
		return runConfig, err
	}
	err = json.Unmarshal([]byte(val), &runConfig)
	return runConfig, err
}

// SetAdminRunConfig will add the session
func (r *RedisRepo) SetAdminRunConfig(key string, runConfig api.RunConfig) (err error) {
	// Marshal to JSON
	var js []byte
	js, err = json.Marshal(runConfig)
	if err != nil {
		return err
	}

	// Set key value
	var OK string
	OK, err = r.redis.Set(r.GetAdminSessionRunKey(key), string(js), adminTTL)
	if err != nil {
		return err
	}
	if OK != "OK" {
		return errors.New("Not OK")
	}
	r.setPort(runConfig.Port, r.GetAdminSessionRunKey(key), userTTL)
	return nil
}

// Functions for adding additional port keys
func (r *RedisRepo) setPort(port string, value string, TTL time.Duration) {
	_, err := r.redis.Set("port:"+port, value, TTL)
	if err != nil {
		//TODO: do not ignore this error.
	}
	log.Printf("[RedisSession]: Added configuration for port: %s\n", port)
}

func (r *RedisRepo) delPort(port string) {
	_, err := r.redis.Del("port:" + port)
	if err != nil {
		//TODO: do not ignore this error.
	}
	log.Printf("[RedisSession]: Removed configuration for port: %s\n", port)
}

// RemoveIncosistentRedisKeys is used when a container is
func (r *RedisRepo) RemoveIncosistentRedisKeys(port string) {
	val, _ := r.redis.Get("port:" + port)
	if val != "" {
		r.delPort(port)
		_, _ = r.redis.Del(val)
	}
}

// ExistsPort is a used by PeriodicChecker function to determine whether a running container should be killed, if the corresponding port key has expired.
func (r *RedisRepo) ExistsPort(port string) bool {
	exists, _ := r.redis.Exists("port:" + port)
	// TODO: YOLO error handling
	return exists
}
