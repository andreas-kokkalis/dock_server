package auth_test

import (
	"flag"
	"net/http"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/auth"
	"github.com/andreas-kokkalis/dock_server/pkg/api/auth/spec/authspec"
	"github.com/andreas-kokkalis/dock_server/pkg/integration"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	dir     string
	testDir string
)

func init() {
	flag.StringVar(&dir, "dir", "../../../../", "dir specifies the relative to the current package position of the top level directory")
}

func TestImageEndpoints(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Image Suite")
}

var (
	spec *integration.Spec
)

var _ = BeforeSuite(func() {
	spec = integration.NewSpec(dir)
	Describe("Initialize config", spec.InitConfig())
	Describe("Connect to postgres", spec.InitDBConnection())
	Describe("Restore database state", spec.RestoreDB())
	Describe("Connect to redis", spec.InitRedisConnection())

	router := httprouter.New()
	authService := auth.NewService(spec.AdminRepo, spec.RedisRepo)
	router.GET("/v0/admin/logout", authService.AdminLogout)
	router.POST("/v0/admin/login", authService.AdminLogin)
	spec.Handler = router
})

var _ = AfterSuite(func() {

})

var _ = Describe("Admin auth endpoints", func() {
	// adminGood := api.Admin{Username: "admin", Password: "kthtest"}
	It("Should login", func() {
		authspec.LoginGoodSpecfunc(spec)
		integration.Report("/admin/login", spec.Time)
	})
	It("Should not error if attempting to login and is already logged in", func() {
		authspec.LoginGoodSpecfunc(spec)
	})
	It("Should logout successfully", func() {
		request := integration.NewRequest(http.MethodGet, "/v0/admin/logout", nil).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusOK, "{}").
			WithSessionCookie("")
		spec.AssertAPICall(request, response)
		integration.Report("/admin/logout", spec.Time)
	})
	It("Should fail to logout if already logged out", func() {
		request := integration.NewRequest(http.MethodGet, "/v0/admin/logout", nil)
		response := integration.NewResponse(http.StatusUnauthorized, authspec.AdminLogoutUnauthorized)
		spec.AssertAPICall(request, response)
	})
	It("Should fail to login with invalid credentials", func() {
		request := integration.NewRequest(http.MethodPost, "/v0/admin/login", api.Admin{Username: "kthtest", Password: "foo"})
		response := integration.NewResponse(http.StatusUnauthorized, authspec.AdminLoginPasswordMismatch)
		spec.AssertAPICall(request, response)
	})
})
