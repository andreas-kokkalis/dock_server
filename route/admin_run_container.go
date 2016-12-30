package route

import (
	"fmt"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
	"github.com/andreas-kokkalis/dock-server/session"
	"github.com/julienschmidt/httprouter"
)

type runRequest struct {
	Username string `json:"user"`
	Password string `json:"pwd"`
}

type runResponse struct {
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ContainerID string `json:"id"`
}

// AdminRunContainer POST
// POST /v0/containers/run
func AdminRunContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Validate imageID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		response.WriteError(res, http.StatusBadRequest, er.InvalidImageID)
		return
	}

	cookie, _ := req.Cookie("ses")
	fmt.Println(cookie.Value)
	sessionExists, err := session.ExistsAdminRunConfig(cookie.Value)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}
	var cfg dc.RunConfig
	username := "admin"
	password := "password"
	if sessionExists {
		cfg, err = session.GetAdminRunConfig(cookie.Value)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, er.ServerError)
			return
		}
		fmt.Printf("exists: %v\n", cfg)
		// Update the TTL
		err = session.SetAdminRunConfig(cookie.Value, cfg)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, er.ServerError)
			return
		}
	} else {
		// Session didn't exist. Admin requested to run a container for the first time.
		// Run container and set the session.

		// Run the container and get the url
		cfg, err = dc.RunContainer(imageID, username, password)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}

		err = session.SetAdminRunConfig(cookie.Value, cfg)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, er.ServerError)
			return
		}
	}
	response.SetData(runResponse{URL: cfg.URL, Username: username, Password: password, ContainerID: cfg.ContainerID})
	res.Write(response.Marshal())
}

//
