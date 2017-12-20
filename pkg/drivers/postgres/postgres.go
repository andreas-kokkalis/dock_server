package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	// postgres dialect
	_ "github.com/lib/pq"
)

const (
	driver = "postgres"
)

// DB ...
type DB struct {
	Conn *sqlx.DB
}

// NewDB ...
func NewDB(connectionString string) (*DB, error) {
	conn, err := sqlx.Open(driver, connectionString)
	if err != nil {
		return &DB{conn}, err
	}
	err = conn.Ping()
	return &DB{conn}, err

}

//nolint
const (
	ErrNamedQuery = "error while executing named query"
	ErrScan       = "error scaning result to data structure"
	ErrRow        = "error occured while iterating over results"
)

var (
	// ErrNoResult is returned when retrieving at least a row was expected but the query returned no results
	ErrNoResult = errors.New("no rows returned when expecting results")
)

// nolint
func (db *DB) QueryRow(query string, val, dest interface{}) error {
	row, err := db.Conn.NamedQuery(query, val)
	if err != nil {
		return errors.Wrap(err, ErrNamedQuery)
	}

	defer row.Close() // nolint: errcheck

	if row.Next() {
		if err = row.StructScan(dest); err != nil {
			return errors.Wrap(err, ErrScan)
		}
	} else {
		return ErrNoResult
	}

	if err = row.Err(); err != nil {
		return errors.Wrap(err, ErrRow)
	}
	return nil
}
