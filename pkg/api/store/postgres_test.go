package store

import (
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/andreas-kokkalis/dock_server/pkg/config"
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

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Mock the SQL connection
	conn, mock, _ := sqlmock.New()
	mock.MatchExpectationsInOrder(true)
	db := &DB{Conn: conn}
	defer func() { _ = db.Conn.Close() }()

	testName := "Query()"
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT (.*) FROM (.*)").WillReturnRows(rows)
	r, err := db.Query("SELECT id FROM test")
	assert.NoError(err, testName)
	testName = "rows.Next()"
	assert.NotNil(r.Next(), testName)
	var id int
	err = r.Scan(&id)
	testName = "rows.Scan()"
	assert.NoError(err, testName)
	assert.Equal(id, 1, testName)
	assert.NoError(mock.ExpectationsWereMet(), testName)
}

func TestQueryRow(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Mock the SQL connection
	conn, mock, _ := sqlmock.New()
	mock.MatchExpectationsInOrder(true)
	db := &DB{Conn: conn}
	defer func() { _ = db.Conn.Close() }()

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Docker")
	mock.ExpectQuery("SELECT (.*) FROM (.*) WHERE id = (.*)").WillReturnRows(rows)
	r := db.QueryRow("SELECT name FROM test WHERE id = $1", 1)
	testName := "rows.Scan()"
	var name string
	err := r.Scan(&name)
	assert.NoError(err, testName)
	assert.Equal("Docker", name, testName)
}
