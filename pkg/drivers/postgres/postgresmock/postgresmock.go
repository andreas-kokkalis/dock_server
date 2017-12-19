package postgresmock

import (
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
	"github.com/jmoiron/sqlx"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// nolint
type MockDB struct {
	Mock sqlmock.Sqlmock
	DB   *postgres.DB
}

// nolint
func NewMockDB() *MockDB {
	conn, mock, _ := sqlmock.New()
	db := &postgres.DB{Conn: sqlx.NewDb(conn, "sqlmock")}
	return &MockDB{mock, db}
}

//nolint
func (m *MockDB) CloseDB() {
	_ = m.DB.Conn.Close()
}
