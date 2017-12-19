package dbutil

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/andreas-kokkalis/dock_server/pkg/db"
	"github.com/pkg/errors"
)

// nolint
const (
	createSchemaScript = "create_schema.sql"
	dropSchemaScript   = "drop_schema.sql"
	insertSchemaScript = "insert_schema.sql"
)

// DBManager models the interaction with the database schema for integration tests.
type DBManager struct {
	ScriptPath string
	DB         *db.DB
}

// NewDBManager initializes a DBManager struct
func NewDBManager(connectionString, topDir string) (*DBManager, error) {
	dbConn, err := db.NewDB(connectionString)
	if err != nil {
		return nil, err
	}
	return &DBManager{
		ScriptPath: topDir,
		DB:         dbConn,
	}, nil
}

func (d *DBManager) loadSQLScript(scriptFile string) (string, error) {

	file := path.Join(d.ScriptPath, scriptFile)

	// Check whether the expected file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "", errors.Wrap(err, "Cannot find file"+file)
	}

	// Load file contents
	script, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(script), nil
}

// CreateSchema creates the database schema
func (d *DBManager) CreateSchema() error {
	createSchema, err := d.loadSQLScript(createSchemaScript)
	if err != nil {
		return errors.Wrap(err, "Unable to read create_schema.sql script")
	}

	tx, err := d.DB.Conn.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(createSchema)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_ = tx.Commit()
	return nil
}

// DropSchema drops the database schema
func (d *DBManager) DropSchema() error {
	dropSchema, err := d.loadSQLScript(dropSchemaScript)
	if err != nil {
		return errors.Wrap(err, "Unable to read drop_schema.sql script")
	}
	tx, err := d.DB.Conn.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(dropSchema)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}

// InsertSchema inserts basic data to the database schema
func (d *DBManager) InsertSchema() error {
	insertSchema, err := d.loadSQLScript(insertSchemaScript)
	if err != nil {
		return errors.Wrap(err, "Unable to read insert_schema.sql script")
	}
	tx, err := d.DB.Conn.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(insertSchema)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_ = tx.Commit()
	return nil
}
