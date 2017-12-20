package store

import (
	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
)

//nolint
type DBAdminRepo struct {
	db *postgres.DB
}

//nolint
func NewDBAdminRepo(db *postgres.DB) *DBAdminRepo {
	return &DBAdminRepo{db}
}

//nolint
func (d *DBAdminRepo) GetAdminByUsername(input api.Admin) (api.Admin, error) {
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
