package route

import (
	"encoding/json"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/julienschmidt/httprouter"
)

// ListImages returns the list of images
func ListImages(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Get the list of images
	images, err := dc.ListImages()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Data = images
	js, err := json.Marshal(images)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Write(js)
}
