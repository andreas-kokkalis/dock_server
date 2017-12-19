package store

import (
	"database/sql"

	"github.com/Davmuz/gqt"
	// postgres dialect
	_ "github.com/lib/pq"
)

const (
	driver = "postgres"
)

// DB ...
type DB struct {
	Conn *sql.DB
}

// NewDB ...
func NewDB(connectionString string) (*DB, error) {
	conn, err := sql.Open(driver, connectionString)
	if err != nil {
		return &DB{conn}, err
	}
	err = conn.Ping()
	_ = gqt.Add("templates/sql", "*.sql")
	return &DB{conn}, err

}

// Query executes an sql query and returns *sql.Rows
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Conn.Query(query, args...)
	return rows, err
}

// QueryRow executes an asql query and returns a single row
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	// TODO: there is a problem with row returned.
	var row *sql.Row
	row = db.Conn.QueryRow(query, args...)
	return row
}
