package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// GetImageHistory returns the history of a particular image
// GET /images/history/abc33412adwq
func GetImageHistory(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Validate imageID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		http.Error(res, er.InvalidImageID, http.StatusUnprocessableEntity)
		response.AddError(er.InvalidImageID)
		res.Write(response.Marshal())
		return
	}

	// Retrieve image history
	history, err := dc.ImageHistory(imageID)
	if err != nil {
		http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
	}

	response.Data = history
	res.Write(response.Marshal())
}
