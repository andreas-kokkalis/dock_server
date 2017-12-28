package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres/postgresmock"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestNewDbAdminRepo(t *testing.T) {
	db := postgresmock.NewMockDB()
	repo := NewAdminDBRepository(db.DB)
	assert.NotNil(t, repo)
}

func TestGetAdminByUsername(t *testing.T) {
	db := postgresmock.NewMockDB()
	repo := NewAdminDBRepository(db.DB)
	db.Mock.ExpectQuery(`SELECT (.+)`).WillReturnError(postgres.ErrNoResult)
	_, err := repo.GetAdminByUsername(api.Admin{Username: "username"})
	assert.Error(t, err)
	assert.NoError(t, db.Mock.ExpectationsWereMet())
	db.CloseDB()

	db = postgresmock.NewMockDB()
	repo = NewAdminDBRepository(db.DB)
	db.Mock.ExpectQuery(`SELECT (.+)`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow("1", "foo"))
	admin, err := repo.GetAdminByUsername(api.Admin{Username: "username"})
	assert.NoError(t, err)
	assert.Equal(t, api.Admin{Username: "username", Password: "foo", ID: 1}, admin)

	assert.NoError(t, db.Mock.ExpectationsWereMet())
	db.CloseDB()
}
