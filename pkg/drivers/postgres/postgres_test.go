package postgres

import (
	"errors"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/andreas-kokkalis/dock_server/pkg/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var validConfigDir = "../../../conf"

func TestNewDB(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	c, _ := config.NewConfig(validConfigDir, "local")

	testName := "Valid driver"
	db, err := NewDB(c.GetPGConnectionString())
	assert.Error(err, testName)
	assert.NotNil(db.Conn, testName)
	assert.Error(db.Conn.Ping(), testName)

}

type mockDB struct {
	mock sqlmock.Sqlmock
	db   *DB
}

// nolint
func NewMockDB() *mockDB {
	conn, mock, _ := sqlmock.New()
	db := &DB{Conn: sqlx.NewDb(conn, "sqlmock")}
	return &mockDB{mock, db}
}

//nolint
func (m *mockDB) CloseDB() {
	_ = m.db.Conn.Close()
}

func TestQueryRowX(t *testing.T) {

	type obj struct {
		Name string `db:"name"`
		ID   int    `db:"id"`
	}
	var dest obj

	query := "SELECT (.*) FROM (.*) WHERE id = (.*)"
	namedQ := "SELECT name FROM test WHERE id = :id"

	// Mock the SQL connection
	m := NewMockDB()
	m.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Docker"))
	err := m.db.QueryRow(namedQ, obj{ID: 1}, &dest)
	assert.NoError(t, err)
	assert.Equal(t, "Docker", dest.Name)
	assert.Nil(t, m.mock.ExpectationsWereMet())
	m.CloseDB()

	m = NewMockDB()
	m.mock.ExpectQuery(query).WillReturnError(errors.New("connection error"))
	err = m.db.QueryRow(namedQ, obj{ID: 1}, &dest)
	assert.Error(t, err, "connection error")
	assert.Nil(t, m.mock.ExpectationsWereMet())
	m.CloseDB()

	m = NewMockDB()
	m.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"blah"}).AddRow("Docker"))
	err = m.db.QueryRow(namedQ, obj{ID: 1}, &dest)
	assert.Error(t, err, "scan error")
	assert.Nil(t, m.mock.ExpectationsWereMet())
	m.CloseDB()

	m = NewMockDB()
	m.mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"blah"}))
	err = m.db.QueryRow(namedQ, obj{ID: 1}, &dest)
	assert.Error(t, err, "no rows")
	assert.Equal(t, err, ErrNoResult)
	assert.Nil(t, m.mock.ExpectationsWereMet())
	m.CloseDB()
}
