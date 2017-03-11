package mock

import "time"

// RedisMock mocks the store.Redis for testing perposes by implementing store.RedisClient
type RedisMock struct {
	GetFn    func(string) (string, error)
	SetFn    func(string, string, time.Duration) (string, error)
	DelFn    func(string) (int64, error)
	ExistsFn func(string) (bool, error)
}

// Get ...
func (r RedisMock) Get(key string) (string, error) {
	return r.GetFn(key)
}

// Set ...
func (r RedisMock) Set(key string, value string, duration time.Duration) (string, error) {
	return r.SetFn(key, value, duration)
}

// Del ...
func (r RedisMock) Del(key string) (int64, error) {
	return r.DelFn(key)
}

// Exists ...
func (r RedisMock) Exists(value string) (bool, error) {
	return r.ExistsFn(value)
}

// NewRedis ...
func NewRedis() RedisMock {
	return RedisMock{}
}

// WithGet ...
func (r RedisMock) WithGet(value string, err error) RedisMock {
	r.GetFn = func(key string) (string, error) {
		return value, err
	}
	return r
}

// WithSet ...
func (r RedisMock) WithSet(value string, err error) RedisMock {
	r.SetFn = func(key string, value string, duration time.Duration) (string, error) {
		return value, err
	}
	return r
}

// WithDel ...
func (r RedisMock) WithDel(value int64, err error) RedisMock {
	r.DelFn = func(key string) (int64, error) {
		return value, err
	}
	return r
}

// WithExists ...
func (r RedisMock) WithExists(value bool, err error) RedisMock {
	r.ExistsFn = func(key string) (bool, error) {
		return value, err
	}
	return r
}
