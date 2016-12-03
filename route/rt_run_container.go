package route

import (
	"encoding/json"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/julienschmidt/httprouter"
)

// RunResult is the result of running a container
type runResponse struct {
	URL string `json:"url"`
}

// XXX: I don't need the name. The name of the repo is the same. Only the tag is needed.

type runRequest struct {
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

	// Validate json request body
	var reqData runRequest
	err := json.NewDecoder(req.Body).Decode(&reqData)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		response.AddError(http.StatusText(http.StatusBadRequest))
		response.SetStatus(http.StatusBadRequest)
		res.Write(response.Marshal())
		return
	}
	// Validate imageID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		response.AddError(er.InvalidImageID)
		response.SetStatus(http.StatusBadRequest)
		res.Write(response.Marshal())
		return
	}
	// Validate Username and Password
	if !vUsername.MatchString(reqData.Username) || !vPassword.MatchString(reqData.Password) {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		response.AddError(er.CredentialsInvalid)
		response.SetStatus(http.StatusBadRequest)
		res.Write(response.Marshal())
		return
	}

	// Run the container and get the url
	cfg, err1 := dc.RunContainer(imageID, reqData.Username, reqData.Password)
	if err1 != nil {
		// http.Error(res, er.ServerError, http.StatusInternalServerError)
		response.AddError(err1.Error())
		response.SetStatus(http.StatusInternalServerError)
		res.Write(response.Marshal())
		return
	}

	response.SetData(runResponse{URL: cfg.URL})
	res.Write(response.Marshal())
}
