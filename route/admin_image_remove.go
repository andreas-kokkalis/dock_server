package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
	"github.com/julienschmidt/httprouter"
)

// RemoveImage removes an image from the registry
// DELETE /images/abc33412adqw
func RemoveImage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Validate imageID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		response.WriteError(res, http.StatusUnprocessableEntity, er.InvalidImageID)
		return
	}

	// Check if there are running containers of that image
	containers, _ := dc.ContainersByImageID(imageID)
	if len(containers) > 0 {
		response.WriteError(res, http.StatusBadRequest, "This image has running containers. Cannot delete it.")
		return
	}

	// Remove Image
	err := dc.RemoveImage(imageID)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}

	res.Write(response.Marshal())
}
