package repositories

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
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/redis"
)

//go:generate moq -out ../repomocks/redis_repo_mock.go -pkg repomocks . RedisRepository

// RedisRepository models the interactions with the cache storage for containers, students and admins
type RedisRepository interface {

	// each redis key has a hardcoded prefix that indicates the purpose of the key.
	StripSessionKeyPrefix(key string) string

	// api for run configuration of students
	UserRunKeyGet(userID string) string
	UserRunConfigDelete(userID string) error
	UserRunConfigExists(userID string) (bool, error)
	UserRunConfigGet(userID string) (runConfig api.RunConfig, err error)
	UserRunConfigSet(userID string, runConfig api.RunConfig) (err error)

	// ui admin session
	AdminSessionKeyCreate(adminID int) string
	AdminSessionExists(key string) (bool, error)
	AdminSessionSet(key string) error
	AdminSessionDelete(key string) error

	// ui admin running containers
	AdminRunConfigExists(key string) (bool, error)
	AdminRunConfigDelete(key string) error
	AdminRunConfigGet(key string) (runConfig api.RunConfig, err error)
	AdminRunConfigSet(key string, runConfig api.RunConfig) (err error)

	// api for portmapper
	PortIsMapped(port string) bool
	DeleteStaleMappedPort(port string)
}

// RedisRepo ...
type RedisRepo struct {
	redis redis.Redis
}

// NewRedisRepo ...
func NewRedisRepo(redis redis.Redis) RedisRepository {
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

// UserRunKeyGet constructs the user key
func (r *RedisRepo) UserRunKeyGet(userID string) string {
	return usrPrefix + userID
}

// UserRunConfigDelete deletes the user session
func (r *RedisRepo) UserRunConfigDelete(userID string) error {

	runConfig, err := r.UserRunConfigGet(r.UserRunKeyGet(userID))
	if err != nil {
		// TODO: parse error
	}
	r.delPort(runConfig.Port)

	_, err = r.redis.Del(r.UserRunKeyGet(userID))
	return err
}

// UserRunConfigExists returns true if there is a session for the particular user
func (r *RedisRepo) UserRunConfigExists(userID string) (bool, error) {
	return r.redis.Exists(r.UserRunKeyGet(userID))
}

// UserRunConfigGet returns the user session
func (r *RedisRepo) UserRunConfigGet(userID string) (runConfig api.RunConfig, err error) {
	var val string
	val, err = r.redis.Get(r.UserRunKeyGet(userID))
	if err != nil {
		return runConfig, err
	}
	err = json.Unmarshal([]byte(val), &runConfig)
	return runConfig, err
}

// UserRunConfigSet will add the session
func (r *RedisRepo) UserRunConfigSet(userID string, runConfig api.RunConfig) (err error) {
	// Marshal to JSON
	var js []byte
	js, _ = json.Marshal(runConfig)

	// Set key value
	var OK string
	OK, err = r.redis.Set(r.UserRunKeyGet(userID), string(js), userTTL)
	if err != nil {
		return err
	}
	if OK != "OK" {
		return errors.New("Not OK")
	}
	r.setPort(runConfig.Port, r.UserRunKeyGet(userID), userTTL)
	return nil
}

/*
	============================
				ADMIN
	============================
*/

// AdminSessionKeyCreate returns the admin session key
func (r *RedisRepo) AdminSessionKeyCreate(adminID int) string {
	h := md5.New()
	_, _ = io.WriteString(h, strconv.Itoa(adminID))
	_, _ = io.WriteString(h, "key")
	s := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(admPrefix + s)
	return admPrefix + s
}

// AdminSessionExists checks if a session exists for that particular adminID
func (r *RedisRepo) AdminSessionExists(key string) (bool, error) {
	return r.redis.Exists(key)
}

// AdminSessionSet will add a key for that admin
func (r *RedisRepo) AdminSessionSet(key string) error {
	_, err := r.redis.Set(key, key, 0)
	return err
}

// AdminSessionDelete will add a key for that admin
func (r *RedisRepo) AdminSessionDelete(key string) error {
	_, err := r.redis.Del(key)
	return err
}

/* ===========================================
	Admin Run Container Session
==============================================*/

// generateAdminRunKey constructs the admin run key
func (r *RedisRepo) generateAdminRunKey(key string) string {
	return admRunPrefix + key
}

// AdminRunConfigDelete deletes the user session
func (r *RedisRepo) AdminRunConfigDelete(key string) error {
	runConfig, _ := r.AdminRunConfigGet(key)
	r.delPort(runConfig.Port)
	_, err := r.redis.Del(r.generateAdminRunKey(key))
	return err
}

// AdminRunConfigExists returns true if there is a session for the particular user
func (r *RedisRepo) AdminRunConfigExists(key string) (bool, error) {
	return r.redis.Exists(r.generateAdminRunKey(key))
}

// AdminRunConfigGet returns the user session
func (r *RedisRepo) AdminRunConfigGet(key string) (runConfig api.RunConfig, err error) {
	var val string
	val, err = r.redis.Get(r.generateAdminRunKey(key))
	if err != nil {
		return runConfig, err
	}
	_ = json.Unmarshal([]byte(val), &runConfig)
	return runConfig, nil
}

// AdminRunConfigSet will add the session
func (r *RedisRepo) AdminRunConfigSet(key string, runConfig api.RunConfig) (err error) {
	// Marshal to JSON
	var js []byte
	js, _ = json.Marshal(runConfig)

	// Set key value
	var OK string
	OK, err = r.redis.Set(r.generateAdminRunKey(key), string(js), adminTTL)
	if err != nil {
		return err
	}
	if OK != "OK" {
		return errors.New("Not OK")
	}
	r.setPort(runConfig.Port, r.generateAdminRunKey(key), userTTL)
	return nil
}

// Functions for adding additional port keys
func (r *RedisRepo) setPort(port string, value string, TTL time.Duration) {
	_, _ = r.redis.Set("port:"+port, value, TTL)
	log.Printf("[RedisSession]: Added configuration for port: %s\n", port)
}

func (r *RedisRepo) delPort(port string) {
	_, _ = r.redis.Del("port:" + port)
	log.Printf("[RedisSession]: Removed configuration for port: %s\n", port)
}

// DeleteStaleMappedPort is used when a container is
func (r *RedisRepo) DeleteStaleMappedPort(port string) {
	val, _ := r.redis.Get("port:" + port)
	if val != "" {
		r.delPort(port)
		_, _ = r.redis.Del(val)
	}
}

// PortIsMapped is a used by PeriodicChecker function to determine whether a running container should be killed, if the corresponding port key has expired.
func (r *RedisRepo) PortIsMapped(port string) bool {
	exists, _ := r.redis.Exists("port:" + port)
	// TODO: YOLO error handling
	return exists
}
