package route

import (
	"encoding/json"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

type commitContainerRequest struct {
	Comment string `json:"comment"`
	Author  string `json:"auth"`
	RefTag  string `json:"tag"`
}

// CommitContainer creates a new image out of a running container
// POST /v0/containers/commit/:id
// JSON data:
//	* Comment
//	* Author
//	* RefTag
func CommitContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Parse containerID
	containerID := params.ByName("id")
	if !vContainerID.MatchString(containerID) {
		response.WriteError(res, http.StatusBadRequest, er.InvalidContainerID)
		return
	}

	// Parse post params
	decoder := json.NewDecoder(req.Body)
	var data commitContainerRequest
	err := decoder.Decode(&data)
	if err != nil {
		response.WriteError(res, http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Validate post params
	if data.Comment == "" || data.Author == "" || data.RefTag == "" {
		response.WriteError(res, http.StatusUnprocessableEntity, er.InvalidPostData)
		return
	}
	// Create the new image
	err = dc.CommitContainer(data.Comment, data.Author, containerID, data.RefTag)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	res.Write(response.Marshal())
}
