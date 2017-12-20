package integration

import (
	"encoding/json"
	"net/http"
	"path"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres/postgresmock"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/redis/redismock"
	"github.com/andreas-kokkalis/dock_server/pkg/util/dbutil"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func newMockManager(m *postgresmock.MockDB) *dbutil.DBManager {
	return &dbutil.DBManager{
		DB:         m.DB,
		ScriptPath: path.Join(topDir, scriptDir),
	}
}

func TestNewSpec(t *testing.T) {
	s := NewSpec(topDir)
	assert.Equal(t, s.TopDir, topDir)
}

var expectedJSONBody = `
{
  "data": {"foo":"bar"}
}
`

func TestAssertAPICall(t *testing.T) {
	RegisterFailHandler(Fail)
	GinkgoRecover()

	s := NewSpec(topDir)
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.WriteHeader(200)

		type SampleData struct {
			Foo string `json:"foo"`
		}

		wb, _ := json.Marshal(api.Response{Data: SampleData{"bar"}})
		w.Write(wb)
	})

	s.Handler = router

	req := NewRequest(http.MethodGet, "/", nil)
	res := NewResponse(200, expectedJSONBody)

	s.AssertAPICall(req, res)
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

			m := postgresmock.NewMockDB()
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

			m := postgresmock.NewMockDB()
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

			s.Redis = redismock.NewRedisMock().WithClose(nil)
			Describe("close redis connection", s.CloseRedisConnection())
		})

	})
})

func TestInitConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}
