package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/julienschmidt/httprouter"
)

// ListImages returns the list of images along with data per image
// GET /v0/images
func ListImages(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Get the list of images
	images, err := dc.ListImages()
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	response.SetData(images)
	res.Write(response.Marshal())
}
