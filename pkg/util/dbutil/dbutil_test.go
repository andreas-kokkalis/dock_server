package dbutil

import (
	"errors"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres/postgresmock"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	confDir   = "../../../conf"
	scriptDir = "../../../scripts/db"
)

func TestNewDbManager(t *testing.T) {
	dbm, err := NewDBManager("foo", scriptDir)
	assert.Error(t, err, "invalid connection string")
	assert.Nil(t, dbm, "dbmanager should be nil")
}

func newMockManager(m *postgresmock.MockDB) *DBManager {
	return &DBManager{
		DB:         m.DB,
		ScriptPath: scriptDir,
	}
}

func TestLoadSQLScript(t *testing.T) {
	db := newMockManager(postgresmock.NewMockDB())

	script, err := db.loadSQLScript(createSchemaScript)
	assert.NoError(t, err, "script exists")
	assert.NotEqual(t, "", script, "script is not empty")

	script, err = db.loadSQLScript("foo")
	assert.Error(t, err, "script does not exist")
	assert.Equal(t, "", script, "script is empty")
}

func TestCreateSchema(t *testing.T) {
	regexQuery := "CREATE TYPE (.+)"

	m := postgresmock.NewMockDB()
	m.Mock.ExpectBegin()
	m.Mock.ExpectExec(regexQuery).WillReturnResult(sqlmock.NewResult(-1, int64(1)))
	m.Mock.ExpectCommit()
	dbm := newMockManager(m)
	assert.NoError(t, dbm.CreateSchema(), "create schema will not error")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	m = postgresmock.NewMockDB()
	m.Mock.ExpectBegin().WillReturnError(errors.New("begin errored"))
	dbm = newMockManager(m)
	assert.Error(t, dbm.CreateSchema(), "begin errors")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	m = postgresmock.NewMockDB()
	m.Mock.ExpectBegin()
	m.Mock.ExpectExec(regexQuery).WillReturnError(errors.New("exec errored"))
	m.Mock.ExpectRollback()
	dbm = newMockManager(m)
	assert.Error(t, dbm.CreateSchema(), "exec errors")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	dbm = newMockManager(m)
	dbm.ScriptPath = "foo"
	assert.Error(t, dbm.CreateSchema(), "Unable to read create schema script")
}

func TestDropSchema(t *testing.T) {
	regexQuery := "DROP SCHEMA public (.+)"

	m := postgresmock.NewMockDB()
	m.Mock.ExpectBegin()
	m.Mock.ExpectExec(regexQuery).WillReturnResult(sqlmock.NewResult(-1, int64(1)))
	m.Mock.ExpectCommit()
	dbm := newMockManager(m)
	assert.NoError(t, dbm.DropSchema(), "drop schema will not error")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	m = postgresmock.NewMockDB()
	m.Mock.ExpectBegin().WillReturnError(errors.New("begin errored"))
	dbm = newMockManager(m)
	assert.Error(t, dbm.DropSchema(), "begin errors")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	m = postgresmock.NewMockDB()
	m.Mock.ExpectBegin()
	m.Mock.ExpectExec(regexQuery).WillReturnError(errors.New("exec errored"))
	m.Mock.ExpectRollback()
	dbm = newMockManager(m)
	assert.Error(t, dbm.DropSchema(), "exec errors")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	dbm = newMockManager(m)
	dbm.ScriptPath = "foo"
	assert.Error(t, dbm.DropSchema(), "Unable to read drop schema script")
}

func TestInsertSchema(t *testing.T) {
	regexQuery := "INSERT INTO admins(.+)"

	m := postgresmock.NewMockDB()
	m.Mock.ExpectBegin()
	m.Mock.ExpectExec(regexQuery).WillReturnResult(sqlmock.NewResult(-1, int64(1)))
	m.Mock.ExpectCommit()
	dbm := newMockManager(m)
	assert.NoError(t, dbm.InsertSchema(), "insert schema will not error")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	m = postgresmock.NewMockDB()
	m.Mock.ExpectBegin().WillReturnError(errors.New("begin errored"))
	dbm = newMockManager(m)
	assert.Error(t, dbm.InsertSchema(), "begin errors")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	m = postgresmock.NewMockDB()
	m.Mock.ExpectBegin()
	m.Mock.ExpectExec(regexQuery).WillReturnError(errors.New("exec errored"))
	m.Mock.ExpectRollback()
	dbm = newMockManager(m)
	assert.Error(t, dbm.InsertSchema(), "exec errors")
	assert.Nil(t, m.Mock.ExpectationsWereMet())
	m.CloseDB()

	dbm = newMockManager(m)
	dbm.ScriptPath = "foo"
	assert.Error(t, dbm.InsertSchema(), "Unable to read insert schema script")
}
