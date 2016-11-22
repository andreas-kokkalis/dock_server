package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// ListImages returns the list of images along with data per image
func ListImages(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Get the list of images
	images, err := dc.ListImages()
	if err != nil {
		http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		res.Write(response.Marshal())
		return
	}

	response.Data = images
	res.Write(response.Marshal())
}
