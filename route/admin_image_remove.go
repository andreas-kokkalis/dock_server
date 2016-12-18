package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// RemoveImage removes an image from the registry
// DELETE /images/abc33412adqw
func RemoveImage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Validate imageID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		response.WriteError(res, http.StatusUnprocessableEntity, er.InvalidImageID)
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
