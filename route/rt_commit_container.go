package route

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/julienschmidt/httprouter"
)

type commitContainerRequest struct {
	Comment     string `json:"com"`
	Author      string `json:"auth"`
	ContainerID string `json:"id"`
	RefTag      string `json:"tag"`
}

// CommitContainer TODO:
func CommitContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	decoder := json.NewDecoder(req.Body)
	var reqData commitContainerRequest
	err := decoder.Decode(&reqData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if reqData.Comment == "" || reqData.Author == "" || reqData.ContainerID == "" || reqData.RefTag == "" {
		http.Error(res, errors.New("Insufficient post arguments").Error(), http.StatusBadRequest)
		return
	}

	err = dc.CommitContainer(reqData.Comment, reqData.Author, reqData.ContainerID, reqData.RefTag)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		response.AddError(err.Error())
	}
	js, err2 := json.Marshal(response)
	if err2 != nil {
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(js)
}
