package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// KillContainer terminates and removes a containerID
// DELETE /v0/containers/kill/abc33412adqw
func KillContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Validate ContainerID
	containerID := params.ByName("id")
	if vContainerID.MatchString(containerID) {
		http.Error(res, er.InvalidContainerID, http.StatusUnprocessableEntity)
		response.AddError(er.InvalidContainerID)
		res.Write(response.Marshal())
		return
	}

	// TODO: consider pushing the following 3 calls to a dc. backend
	err := dc.StopContainer(containerID)
	if err != nil {
		http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		res.Write(response.Marshal())
		return
	}
	err = dc.KillContainer(containerID)
	if err != nil {
		http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		res.Write(response.Marshal())
		return
	}
	err = dc.RemoveContainer(containerID)
	if err != nil {
		http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		res.Write(response.Marshal())
		return
	}

	response.Data = "OK"
	res.Write(response.Marshal())
}
