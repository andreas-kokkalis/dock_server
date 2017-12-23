package repomocks

import "github.com/andreas-kokkalis/dock_server/pkg/api"

//nolint
func NewRedisRepositoryMock() *RedisRepositoryMock {
	return &RedisRepositoryMock{}
}

//nolint
func (r *RedisRepositoryMock) WithStripSessionKeyPrefix(out string) *RedisRepositoryMock {
	r.StripSessionKeyPrefixFunc = func(_ string) string {
		return out
	}
	return r
}
func (r *RedisRepositoryMock) WithUserRunKeyGet(out string) *RedisRepositoryMock {
	r.UserRunKeyGetFunc = func(_ string) string {
		return out
	}
	return r
}
func (r *RedisRepositoryMock) WithUserRunConfigDelete(err error) *RedisRepositoryMock {
	r.UserRunConfigDeleteFunc = func(_ string) error {
		return err
	}
	return r
}
func (r *RedisRepositoryMock) WithUserRunConfigExists(exists bool, err error) *RedisRepositoryMock {
	r.UserRunConfigExistsFunc = func(_ string) (bool, error) {
		return exists, err
	}
	return r
}
func (r *RedisRepositoryMock) WithUserRunConfigGet(cfg api.RunConfig, err error) *RedisRepositoryMock {
	r.UserRunConfigGetFunc = func(_ string) (api.RunConfig, error) {
		return cfg, err
	}
	return r
}
func (r *RedisRepositoryMock) WithUserRunConfigSet(err error) *RedisRepositoryMock {
	r.UserRunConfigSetFunc = func(_ string, _ api.RunConfig) error {
		return err
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminSessionKeyCreate(str string) *RedisRepositoryMock {
	r.AdminSessionKeyCreateFunc = func(_ int) string {
		return str
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminSessionExists(exists bool, err error) *RedisRepositoryMock {
	r.AdminSessionExistsFunc = func(_ string) (bool, error) {
		return exists, err
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminSessionSet(err error) *RedisRepositoryMock {
	r.AdminSessionSetFunc = func(_ string) error {
		return err
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminSessionDelete(err error) *RedisRepositoryMock {
	r.AdminSessionDeleteFunc = func(_ string) error {
		return err
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminRunConfigExists(exists bool, err error) *RedisRepositoryMock {
	r.AdminRunConfigExistsFunc = func(_ string) (bool, error) {
		return exists, err
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminRunConigDelete(err error) *RedisRepositoryMock {
	r.AdminRunConfigDeleteFunc = func(_ string) error {
		return err
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminRunConfigGet(cfg api.RunConfig, err error) *RedisRepositoryMock {
	r.AdminRunConfigGetFunc = func(_ string) (api.RunConfig, error) {
		return cfg, err
	}
	return r
}
func (r *RedisRepositoryMock) WithAdminRunConfigSet(err error) *RedisRepositoryMock {
	r.AdminRunConfigSetFunc = func(_ string, _ api.RunConfig) error {
		return err
	}
	return r
}
func (r *RedisRepositoryMock) WithPortIsMapped(isMapped bool) *RedisRepositoryMock {
	r.PortIsMappedFunc = func(_ string) bool {
		return isMapped
	}
	return r
}
func (r *RedisRepositoryMock) WithDeleteStaleMappedPort() *RedisRepositoryMock {
	r.DeleteStaleMappedPortFunc = func(_ string) {}
	return r
}
