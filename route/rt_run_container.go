package route

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/andreas-kokkalis/dock-server/dc"
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

// XXX: validator for imageID
// TODO: might need to extend this with
//	* session data.
//	* username for user
//	* password for user
// 		If a session already exists, then you shouldn't run
// 		a new container, but recreate the url for the old container

var validPassword = regexp.MustCompile(`^([a-zA-Z0-9]){5,6}$`)
var validUsername = regexp.MustCompile(`^([a-zA-Z0-9]){2,16}$`)

// RunContainer POST
func RunContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	decoder := json.NewDecoder(req.Body)
	var reqData runRequest
	err := decoder.Decode(&reqData)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if reqData.Name == "" || reqData.RefTag == "" || reqData.Username == "" || reqData.Password == "" ||
		validUsername.MatchString(reqData.Username) == false || validPassword.MatchString(reqData.Password) == false {
		http.Error(res, errors.New("Insufficient post arguments").Error(), http.StatusBadRequest)
		return
	}

	var url string
	url, err = dc.RunContainer(reqData.Name, reqData.RefTag, reqData.Username, reqData.Password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		response.AddError(err.Error())
	}
	response.Data = runResponse{URL: url}
	js, err2 := json.Marshal(response)
	if err2 != nil {
		http.Error(res, err2.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(js)
}
