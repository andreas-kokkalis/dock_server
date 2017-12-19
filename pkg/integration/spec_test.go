package integration

import (
	"path"
	"testing"

	"github.com/andreas-kokkalis/dock_server/cmd/dock_server/schema/dbutil"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type MockDB struct {
	Mock sqlmock.Sqlmock
	DB   *store.DB
}

func newMockDB() *MockDB {
	conn, mock, _ := sqlmock.New()
	db := &store.DB{Conn: conn}
	return &MockDB{mock, db}
}

func (m *MockDB) CloseDB() {
	_ = m.DB.Conn.Close()
}

func newMockManager(m *MockDB) *dbutil.DBManager {
	return &dbutil.DBManager{
		DB:         m.DB,
		ScriptPath: path.Join(topDir, scriptDir),
	}
}

func TestNewSpec(t *testing.T) {
	s := NewSpec(topDir)
	assert.Equal(t, s.TopDir, topDir)
}

var _ = Describe("Test methods of Spec struct", func() {
	Context("Testing spec functions", func() {
		It("Should initialize configuration correctly", func() {
			s := NewSpec(topDir)
			Describe("init config", s.InitConfig())
			Expect(s.Config.GetAPIServerPort()).To(Equal(":8080"))
		})
		It("Should close the database connection", func() {
			s := NewSpec(topDir)
			Describe("init config", s.InitConfig())

			m := newMockDB()
			dbm := newMockManager(m)
			s.DBManager = dbm

			m.Mock.ExpectClose()
			Describe("close db connection", s.CloseDBConnection())
			Expect(m.Mock.ExpectationsWereMet()).To(BeNil())
			m.CloseDB()
		})
		It("Should restore the database", func() {
			s := NewSpec(topDir)
			Describe("init config", s.InitConfig())

			m := newMockDB()
			dbm := newMockManager(m)
			s.DBManager = dbm

			m.Mock.ExpectBegin()
			m.Mock.ExpectExec("DROP SCHEMA public (.+)").WillReturnResult(sqlmock.NewResult(-1, int64(1)))
			m.Mock.ExpectCommit()
			m.Mock.ExpectBegin()
			m.Mock.ExpectExec("CREATE TYPE (.+)").WillReturnResult(sqlmock.NewResult(-1, int64(1)))
			m.Mock.ExpectCommit()
			m.Mock.ExpectBegin()
			m.Mock.ExpectExec("INSERT INTO admins(.+)").WillReturnResult(sqlmock.NewResult(-1, int64(1)))
			m.Mock.ExpectCommit()
			Describe("restore Database", s.RestoreDB())
			Expect(m.Mock.ExpectationsWereMet()).To(BeNil())
			m.CloseDB()
		})
		It("Should close the redis connection", func() {
			s := NewSpec(topDir)
			Describe("init config", s.InitConfig())

			m := mock.NewRedis()
			m.CloseFn = func() error { return nil }
			s.Redis = m
			Describe("close redis connection", s.CloseRedisConnection())
		})

	})
})

func TestInitConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
