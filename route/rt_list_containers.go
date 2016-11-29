package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// GetContainers returns list of containers by status.
// GET /v0/containers
// GET /v0/containers/:status
func GetContainers(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Validate if status has any of the accepted input
	status := params.ByName("status")
	if vContainerState.MatchString(status) == false {
		// http.Error(res, er.InvalidContainerState, http.StatusUnprocessableEntity)
		response.AddError(er.InvalidContainerState)
		response.SetStatus(http.StatusUnprocessableEntity)
		res.Write(response.Marshal())
		return
	}

	// Get the list of containers
	containers, err := dc.GetContainers(status)
	if err != nil {
		// http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		response.SetStatus(http.StatusInternalServerError)
		res.Write(response.Marshal())
		return
	}

	response.SetData(containers)
	res.Write(response.Marshal())
}
