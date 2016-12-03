package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// CreateRunURL will create a LTI basic launch URL and return the configuration
func CreateRunURL(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Validate ContainerID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		response.AddError(er.InvalidImageID)
		response.SetStatus(http.StatusBadRequest)
		res.Write(response.Marshal())
		return
	}

	tag, err := dc.GetTagByID(imageID)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		response.AddError(er.ServerError)
		response.SetStatus(http.StatusInternalServerError)
		res.Write(response.Marshal())
	}
	if tag == "" {
		response.AddError("Image id does not exist")
		res.Write(response.Marshal())
	} else {
		response.SetData("https://localhost:8080/lti/launch/" + imageID)
	}
	res.Write(response.Marshal())
}
