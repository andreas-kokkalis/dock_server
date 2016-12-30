package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
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
		response.WriteError(res, http.StatusBadRequest, er.InvalidImageID)
		return
	}

	tag, err := dc.GetTagByID(imageID)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	if tag == "" {
		response.WriteError(res, http.StatusFailedDependency, er.ImageNotFound)
		return
	}
	response.SetData("https://localhost:8080/lti/launch/" + imageID)
	res.Write(response.Marshal())
}
