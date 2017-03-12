package image

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/pkg/api"
	"github.com/andreas-kokkalis/dock-server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock-server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/route/er"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	db     *store.DB
	redis  *store.RedisRepo
	docker *docker.Repo
}

// NewService creates a new Image Service
func NewService(db *store.DB, redis *store.RedisRepo, docker *docker.Repo) Service {
	return Service{db, redis, docker}
}

// ListImages returns the list of images along with data per image
// GET /v0/images
func ListImages(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// Get the list of images
		images, err := s.docker.ImageList()
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}
		response.SetData(images)
		res.Write(response.Marshal())
	}
}

// GetImageHistory returns the history of a particular image
// GET /images/history/:id
func GetImageHistory(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// TODO: auth

		// Validate imageID
		imageID := params.ByName("id")
		if !api.VImageID.MatchString(imageID) {
			response.WriteError(res, http.StatusBadRequest, api.ErrInvalidImageID)
			return
		}

		// Retrieve image history
		history, err := s.docker.ImageHistory(imageID)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}
		response.SetData(history)
		res.Write(response.Marshal())
	}
}

// RemoveImage removes an image from the registry
// DELETE /images/abc33412adqw
func RemoveImage(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// Validate imageID
		imageID := params.ByName("id")
		if !api.VImageID.MatchString(imageID) {
			response.WriteError(res, http.StatusUnprocessableEntity, er.InvalidImageID)
			return
		}

		// Check if there are running containers of that image
		containers, _ := s.docker.ContainersByImageID(imageID)
		if len(containers) > 0 {
			response.WriteError(res, http.StatusBadRequest, "This image has running containers. Cannot delete it.")
			return
		}

		// Remove Image
		err := s.docker.RemoveImage(imageID)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}

		res.Write(response.Marshal())
	}
}

/*
// CreateRunURL will create a LTI basic launch URL and return the configuration
func CreateRunURL(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// TODO: auth

		// Validate ContainerID
		imageID := params.ByName("id")
		if !api.VImageID.MatchString(imageID) {
			response.WriteError(res, http.StatusBadRequest, api.ErrInvalidImageID)
			return
		}

		tag, err := s.docker.GetTagByID(imageID)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}
		if tag == "" {
			response.WriteError(res, http.StatusFailedDependency, er.ImageNotFound)
			return
		}
		response.SetData("https://localhost:8080/lti/launch/" + imageID)
		res.Write(response.Marshal())
	}
}
*/
