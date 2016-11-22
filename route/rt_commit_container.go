package route

import (
	"encoding/json"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

type commitContainerRequest struct {
	Comment     string `json:"com"`
	Author      string `json:"auth"`
	ContainerID string `json:"id"`
	RefTag      string `json:"tag"`
}

// CommitContainer creates a new image out of a running container
// POST containers/commit
// data:
//	* Comment
//	* Author
//	* ContainerID
//	* RefTag
func CommitContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Parse post params
	decoder := json.NewDecoder(req.Body)
	var reqData commitContainerRequest
	err := decoder.Decode(&reqData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Validate post params
	if reqData.Comment == "" || reqData.Author == "" || reqData.ContainerID == "" || reqData.RefTag == "" {
		http.Error(res, er.InvalidPostData, http.StatusUnprocessableEntity)
		response.AddError(er.InvalidPostData)
		res.Write(response.Marshal())
		return
	}

	// Create the new image
	err = dc.CommitContainer(reqData.Comment, reqData.Author, reqData.ContainerID, reqData.RefTag)
	if err != nil {
		http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err.Error())
		res.Write(response.Marshal())
		return
	}

	response.Data = "OK"
	res.Write(response.Marshal())
}
