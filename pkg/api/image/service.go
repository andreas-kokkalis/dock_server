package image

import (
	"net/http"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	redis  repositories.RedisRepository
	docker repositories.DockerRepository
}

// NewService creates a new Image Service
func NewService(redis repositories.RedisRepository, docker repositories.DockerRepository) Service {
	return Service{redis, docker}
}

// ListImages returns the list of images along with data per image
// GET /v0/images
func (s Service) ListImages(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	images, err := s.docker.ImageList()
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.WriteOKResponse(w, images)
}

// GetImageHistory returns the history of a particular image
// GET /images/history/:id
func (s Service) GetImageHistory(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Validate imageID
	imageID := params.ByName("id")
	if !api.VImageID.MatchString(imageID) {
		api.WriteErrorResponse(w, http.StatusBadRequest, api.ErrInvalidImageID)
		return
	}

	// Retrieve image history
	history, err := s.docker.ImageHistory(imageID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.WriteOKResponse(w, history)
}

var (
	// ErrImageHasRunningContainers is returned when attempting to delete an image that is in use.
	ErrImageHasRunningContainers = "Cannot delete an image that is currently used by running containers."
)

// RemoveImage removes an image from the registry
// DELETE /images/abc33412adqw
func (s Service) RemoveImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Validate imageID
	imageID := params.ByName("id")
	if !api.VImageID.MatchString(imageID) {
		api.WriteErrorResponse(w, http.StatusBadRequest, api.ErrInvalidImageID)
		return
	}

	// Check if there are running containers of that image
	containers, _ := s.docker.GetRunningContainersByImageID(imageID)
	if len(containers) > 0 {
		api.WriteErrorResponse(w, http.StatusBadRequest, ErrImageHasRunningContainers)
		return
	}

	// Remove Image
	err := s.docker.ImageRemove(imageID)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.WriteOKResponse(w, "Image was removed successfuly")
}
