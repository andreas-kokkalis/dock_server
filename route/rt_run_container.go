package route

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// RunResult is the result of running a container
type runResponse struct {
	URL string `json:"url"`
}

type runRequest struct {
	Name     string `json:"name"`
	RefTag   string `json:"tag"`
	Username string `json:"user"`
	Password string `json:"pwd"`
}

// TODO: might need to extend this with
//	* session data.
//	* username for user
//	* password for user
// 		If a session already exists, then you shouldn't run
// 		a new container, but recreate the url for the old container

// RunContainer POST
// POST /v0/containers/run
func RunContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	decoder := json.NewDecoder(req.Body)
	var reqData runRequest
	err := decoder.Decode(&reqData)
	if err != nil {
		http.Error(res, er.InvalidPostData, http.StatusUnprocessableEntity)
		return
	}

	// XXX: might not need to get the name as a parameter
	if reqData.Name == "" || reqData.RefTag == "" ||
		vUsername.MatchString(reqData.Username) == false || vPassword.MatchString(reqData.Password) == false {
		http.Error(res, errors.New("Insufficient post arguments").Error(), http.StatusBadRequest)
		return
	}

	// Run the container and get the url
	url, err1 := dc.RunContainer(reqData.Name, reqData.RefTag, reqData.Username, reqData.Password)
	if err1 != nil {
		http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err1.Error())
		res.Write(response.Marshal())
		return
	}

	response.Data = runResponse{URL: url}
	res.Write(response.Marshal())
}
