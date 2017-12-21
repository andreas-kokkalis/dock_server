package repomocks

import "github.com/andreas-kokkalis/dock_server/pkg/api"

// NewAdminDBRepositoryMock initializes a mock AdminDBRepository
func NewAdminDBRepositoryMock() *AdminDBRepositoryMock {
	return &AdminDBRepositoryMock{}
}

// WithGetAdminByUsername sets the GetAdminByUsername mock function
func (a *AdminDBRepositoryMock) WithGetAdminByUsername(admin api.Admin, err error) *AdminDBRepositoryMock {
	a.GetAdminByUsernameFunc = func(_ api.Admin) (api.Admin, error) {
		return admin, err
	}
	return a
}
