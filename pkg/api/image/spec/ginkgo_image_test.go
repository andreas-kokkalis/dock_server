package image_test

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/auth"
	"github.com/andreas-kokkalis/dock_server/pkg/api/auth/spec/authspec"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image/spec/imgspec"
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
	Describe("Init docker repo", spec.InitDockerRepo())

	router := httprouter.New()
	authService := auth.NewService(spec.AdminRepo, spec.RedisRepo)
	// Image service
	imageService := image.NewService(spec.RedisRepo, spec.DockerRepo)

	router.POST("/v0/admin/login", authService.AdminLogin)
	router.GET("/v0/admin/images", authService.SessionAuth(imageService.ListImages))
	router.GET("/v0/admin/images/history/:id", authService.SessionAuth(imageService.GetImageHistory))
	router.DELETE("/v0/admin/images/delete/:id", authService.SessionAuth(imageService.RemoveImage))

	spec.Handler = router
})

var _ = AfterSuite(func() {

})

var _ = Describe("Image api endpoints", func() {
	It("Should fail to get any context without a valid session", func() {
		request := integration.NewRequest(http.MethodGet, "/v0/admin/images", nil)
		response := integration.NewResponse(http.StatusUnauthorized, authspec.AdminLoginUnauthorized)
		spec.AssertAPICall(request, response)
	})
	It("Should login", func() {
		authspec.LoginGoodSpecfunc(spec)
	})

	var img api.Img
	It("Should list all images", func() {
		request := integration.NewRequest(http.MethodGet, "/v0/admin/images", nil).WithSessionCookie(authspec.ValidSessionKey)
		response := integration.NewResponse(http.StatusOK, imgspec.ImageListGood)
		spec.AssertAPICall(request, response)

		var images []api.Img
		response.Unmarshall(&images)
		img = images[len(images)-1]

	})
	It("Should get image history of seed image", func() {
		request := integration.NewRequest(http.MethodGet, fmt.Sprintf("/v0/admin/images/history/%s", img.ID), nil).WithSessionCookie(authspec.ValidSessionKey)
		response := integration.NewResponse(http.StatusOK, imgspec.ImageHistoryGood)
		spec.AssertAPICall(request, response)
	})
	It("Should not find image history for invalid image ID", func() {
		invalidImageIDReponse := `
		{
	  		"errors": [
	    		"ImageID is invalid"
	  		],
	  		"status": "Bad Request"
		}`
		request := integration.NewRequest(http.MethodGet, "/v0/admin/images/history/123", nil).WithSessionCookie(authspec.ValidSessionKey)
		response := integration.NewResponse(http.StatusBadRequest, invalidImageIDReponse)
		spec.AssertAPICall(request, response)
	})
})
