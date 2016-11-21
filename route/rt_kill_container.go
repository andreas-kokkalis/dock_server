package route

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/julienschmidt/httprouter"
)

// KillContainer route
func KillContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	containerID := params.ByName("id")
	if containerID == "" {
		http.Error(res, errors.New("Invalid ContainerID").Error(), http.StatusBadRequest)
		return
	}

	err := dc.StopContainer(containerID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err = dc.KillContainer(containerID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Data = "Container killed successfully"
	js, err2 := json.Marshal(response)
	if err2 != nil {
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(js)
}
