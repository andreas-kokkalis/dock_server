package route

import (
	"encoding/json"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

type commitContainerRequest struct {
	Comment string `json:"com"`
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

	// Parse post params
	containerID := params.ByName("id")
	decoder := json.NewDecoder(req.Body)
	var reqData commitContainerRequest
	err := decoder.Decode(&reqData)
	// Validate post params
	if err != nil || reqData.Comment == "" || reqData.Author == "" || !vContainerID.MatchString(containerID) || reqData.RefTag == "" {
		// http.Error(res, er.InvalidPostData, http.StatusUnprocessableEntity)
		response.AddError(er.InvalidPostData)
		response.SetStatus(http.StatusUnprocessableEntity)
		res.Write(response.Marshal())
		return
	}

	// Create the new image
	err = dc.CommitContainer(reqData.Comment, reqData.Author, containerID, reqData.RefTag)
	if err != nil {
		// http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		response.SetStatus(http.StatusInternalServerError)
		res.Write(response.Marshal())
		return
	}

	response.SetStatus(http.StatusOK)
	res.Write(response.Marshal())
}
