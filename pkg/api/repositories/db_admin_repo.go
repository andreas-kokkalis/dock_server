package repositories

import (
	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
)

//go:generate moq -out ../repomocks/db_admin_repo_mock.go -pkg repomocks . AdminDBRepository

// AdminDBRepository models the interaction of an admin with the database
type AdminDBRepository interface {
	GetAdminByUsername(input api.Admin) (api.Admin, error)
}

// AdminDBRepo implements AdminDBRepository
type AdminDBRepo struct {
	db *postgres.DB
}

// NewAdminDBRepository initializes an AdminDBRepo
func NewAdminDBRepository(db *postgres.DB) AdminDBRepository {
	return &AdminDBRepo{db}
}

// GetAdminByUsername returnes an admin object if the username matches with one registered in the database.
func (d *AdminDBRepo) GetAdminByUsername(input api.Admin) (api.Admin, error) {
	q := `
		SELECT
			id,
			password
		FROM
			admins
		WHERE
			username = :username
		`
	var admin api.Admin
	err := d.db.QueryRow(q, input, &admin)
	if err != nil && err != postgres.ErrNoResult {
		return api.Admin{}, err
	}

	return admin, nil
}
