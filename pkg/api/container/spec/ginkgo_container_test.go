package container_test

import (
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/container"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image"
	"github.com/andreas-kokkalis/dock_server/pkg/api/image/spec/imgspec"
	"github.com/andreas-kokkalis/dock_server/pkg/api/portmapper"
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
	Describe("Init port mapper", spec.InitPortMapper())

	router := httprouter.New()
	// Image service
	imageService := image.NewService(spec.RedisRepo, spec.DockerRepo)
	router.GET("/v0/admin/images", imageService.ListImages)
	router.GET("/v0/admin/images/history/:id", imageService.GetImageHistory)
	router.DELETE("/v0/admin/images/delete/:id", imageService.RemoveImage)

	cntService := container.NewService(spec.RedisRepo, spec.DockerRepo, spec.Mapper)
	router.POST("/v0/admin/containers/run/:id", cntService.AdminRunContainer)
	router.DELETE("/v0/admin/containers/kill/:id", cntService.AdminKillContainer)
	router.POST("/v0/admin/containers/commit/:id", cntService.CommitContainer)
	router.GET("/v0/admin/containers/list", cntService.GetContainers)
	router.GET("/v0/admin/containers/list/:status", cntService.GetContainers)

	spec.Handler = router
})

var _ = AfterSuite(func() {

})

var _ = Describe("Image Suite", func() {

	BeforeEach(func() {
		portmapper.Check(spec.DockerRepo, spec.Mapper, spec.RedisRepo)
	})
	AfterEach(func() {
		portmapper.Check(spec.DockerRepo, spec.Mapper, spec.RedisRepo)
	})

	var img api.Img
	It("Should list all images", func() {
		request := integration.NewRequest(http.MethodGet, "/v0/admin/images", nil)
		response := integration.NewResponse(http.StatusOK, imgspec.ImageListGood)
		spec.AssertAPICall(request, response)

		var images []api.Img
		response.Unmarshall(&images)
		img = images[len(images)-1]

	})
	It("Should fail running container of invalid image ID", func() {
		failContainerRunRequest := `
		{
	      "errors": [
	        "ImageID is invalid"
	      ],
	      "status": "Bad Request"
	    }`
		request := integration.NewRequest(http.MethodPost, fmt.Sprintf("/v0/admin/containers/run/%s", "foobar"), nil)
		response := integration.NewResponse(http.StatusBadRequest, failContainerRunRequest)
		spec.AssertAPICall(request, response)
	})
	It("Should fail running container without a valid session cookie", func() {
		failContainerRunRequest := `
		{
	      "errors": [
	        "admin session cookie is not set"
	      ],
	      "status": "Bad Request"
	    }`
		request := integration.NewRequest(http.MethodPost, fmt.Sprintf("/v0/admin/containers/run/%s", img.ID), nil)
		response := integration.NewResponse(http.StatusBadRequest, failContainerRunRequest)
		spec.AssertAPICall(request, response)
	})
	var containerRun api.ContainerRun
	var containerRunResponse = `
	{
      "data": {
        "id": "([A-Fa-f0-9]{12,64})$",
        "password": "password",
        "url": "https://127.0.0.1:([0-9]{4})$",
        "username": "guest"
      }
    }`
	It("Should run a container", func() {
		request := integration.NewRequest(http.MethodPost, fmt.Sprintf("/v0/admin/containers/run/%s", img.ID), nil).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusOK, containerRunResponse)
		spec.AssertAPICall(request, response)
		response.Unmarshall(&containerRun)
	})
	var containerKillResponse = `
	{
      "data": "Container Killed"
    }`
	It("Should remove container", func() {
		request := integration.NewRequest(http.MethodDelete, fmt.Sprintf("/v0/admin/containers/kill/%s", containerRun.ContainerID), nil).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusOK, containerKillResponse)
		spec.AssertAPICall(request, response)
	})
	It("Should run the container again after removing it", func() {
		request := integration.NewRequest(http.MethodPost, fmt.Sprintf("/v0/admin/containers/run/%s", img.ID), nil).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusOK, containerRunResponse)
		spec.AssertAPICall(request, response)
		response.Unmarshall(&containerRun)
	})
	It("Should fail committing the container with invalid JSON body", func() {
		commitFailResponse := `
		{
	      "errors": [
	        "POST parameters have errors"
	      ],
	      "status": "Bad Request"
	    }`
		request := integration.NewRequest(http.MethodPost, fmt.Sprintf("/v0/admin/containers/commit/%s", containerRun.ContainerID), nil).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusBadRequest, commitFailResponse)
		spec.AssertAPICall(request, response)
	})
	//TODO Run the container again
	var containerCommitResponse = `
	{
      "data": {
        "imageID": "([A-Fa-f0-9]{12,64})$"
      }
    }`
	var newImgID api.ContainerCommitResponse
	It("Should commit the container", func() {
		request := integration.NewRequest(
			http.MethodPost,
			fmt.Sprintf("/v0/admin/containers/commit/%s", containerRun.ContainerID), api.ContainerCommitRequest{Comment: "comment", Author: "Author", RefTag: "testref"}).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusOK, containerCommitResponse)
		spec.AssertAPICall(request, response)
		response.Unmarshall(&newImgID)
	})
	It("Should fail removing the container since commit removed it already", func() {
		failContainerRemoveRequest := `
		{
	      "errors": [
	        "Container did not exist."
	      ],
	      "status": "Internal Server Error"
	    }`
		request := integration.NewRequest(http.MethodDelete, fmt.Sprintf("/v0/admin/containers/kill/%s", containerRun.ContainerID), nil).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusInternalServerError, failContainerRemoveRequest)
		spec.AssertAPICall(request, response)
	})
	It("Should remove the new image", func() {
		imageRemoveResponse := `
		{
	      "data": "Image was removed successfuly"
	    }`
		request := integration.NewRequest(http.MethodDelete, fmt.Sprintf("/v0/admin/images/delete/%s", newImgID.ImageID), nil).
			WithSessionCookie("adm:7ff10abb653dead4186089acbd2b7891")
		response := integration.NewResponse(http.StatusOK, imageRemoveResponse)
		spec.AssertAPICall(request, response)
	})
})
