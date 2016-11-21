package route

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/julienschmidt/httprouter"
)

// RemoveImage removes an image from the registry
func RemoveImage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	imageID := params.ByName("id")
	if imageID == "" {
		http.Error(res, errors.New("Insufficient post arguments").Error(), http.StatusBadRequest)
		return
	}

	err := dc.RemoveImage(imageID)
	if err != nil {
		response.AddError(err.Error())
	}

	js, err2 := json.Marshal(response)
	if err2 != nil {
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(js)
}
