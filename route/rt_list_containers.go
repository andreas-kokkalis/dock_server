package route

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/julienschmidt/httprouter"
)

// Match only against empty string or any of the given values
var validStatus = regexp.MustCompile(`^(|created|restarting|running|paused|exited|dead)$`)

// GetContainers returns list of running containers
func GetContainers(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Validate if status has any of the accepted input
	status := params.ByName("status")
	if validStatus.MatchString(status) == false {
		response.AddError("Status does not match the required input")

		js, err := json.Marshal(response)
		if err != nil {
			// http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Write(js)
		return
	}

	// Get the list of containers
	containers, err := dc.GetContainers(status)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	response.Data = containers

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(js)
}
