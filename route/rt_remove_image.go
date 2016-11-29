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
		// http.Error(res, er.InvalidImageID, http.StatusUnprocessableEntity)
		response.AddError(er.InvalidImageID)
		response.SetStatus(http.StatusUnprocessableEntity)
		res.Write(response.Marshal())
		return
	}

	// Remove Image
	err := dc.RemoveImage(imageID)
	if err != nil {
		// http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		response.SetStatus(http.StatusInternalServerError)
		res.Write(response.Marshal())
		return
	}

	res.Write(response.Marshal())
}
