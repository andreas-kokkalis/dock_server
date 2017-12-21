package image_test

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image"
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
	// TODO: remove these on once used in the proper testing environment
	// Describe("Connect to postgres", spec.InitDBConnection())
	// Describe("Restore database state", spec.RestoreDB())
	Describe("Connect to redis", spec.InitRedisConnection())
	Describe("Init docker repo", spec.InitDockerRepo())

	router := httprouter.New()
	// Image service
	imageService := image.NewService(spec.RedisRepo, spec.DockerRepo)
	router.GET("/v0/admin/images", image.ListImages(imageService))
	router.GET("/v0/admin/images/history/:id", image.GetImageHistory(imageService))
	router.DELETE("/v0/admin/images/delete/:id", image.RemoveImage(imageService))

	spec.Handler = router
})

var _ = AfterSuite(func() {

})

var _ = Describe("Image Suite", func() {
	var img *api.Img
	It("Should list images", func() {
		request := integration.NewRequest(http.MethodGet, "/v0/admin/images", nil)
		response := integration.NewResponse(http.StatusOK, imageListGood)
		spec.AssertAPICall(request, response)

		type res struct {
			Data []api.Img
		}
		var result res
		response.ToStructure(&result)
		img = &result.Data[len(result.Data)-1]
	})
	It("Should get image history", func() {
		request := integration.NewRequest(http.MethodGet, fmt.Sprintf("/v0/admin/images/history/%s", img.ID), nil)
		response := integration.NewResponse(http.StatusOK, imageHistoryGood)
		spec.AssertAPICall(request, response)
	})
})
