package redismock

import "time"

// NewRedisMock initializes an empty RedisMock
func NewRedisMock() *RedisMock {
	return &RedisMock{}
}

// WithGet ...
func (r *RedisMock) WithGet(value string, err error) *RedisMock {
	r.GetFunc = func(_ string) (string, error) {
		return value, err
	}
	return r
}

// WithSet ...
func (r *RedisMock) WithSet(response string, err error) *RedisMock {
	r.SetFunc = func(_ string, _ string, _ time.Duration) (string, error) {
		return response, err
	}
	return r
}

// WithDel ...
func (r *RedisMock) WithDel(status int64, err error) *RedisMock {
	r.DelFunc = func(_ string) (int64, error) {
		return status, err
	}
	return r
}

// WithExists ...
func (r *RedisMock) WithExists(value bool, err error) *RedisMock {
	r.ExistsFunc = func(_ string) (bool, error) {
		return value, err
	}
	return r
}

// WithClose ...
func (r *RedisMock) WithClose(err error) *RedisMock {
	r.CloseFunc = func() error {
		return err
	}
	return r
}
