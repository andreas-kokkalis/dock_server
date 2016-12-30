package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
	"github.com/julienschmidt/httprouter"
)

// GetImageHistory returns the history of a particular image
// GET /images/history/:id
func GetImageHistory(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Validate imageID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		response.WriteError(res, http.StatusBadRequest, er.InvalidImageID)
		return
	}

	// Retrieve image history
	history, err := dc.ImageHistory(imageID)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	response.SetData(history)
	res.Write(response.Marshal())
}
