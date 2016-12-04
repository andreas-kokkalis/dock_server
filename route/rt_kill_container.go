package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/andreas-kokkalis/dock-server/session"
	"github.com/julienschmidt/httprouter"
)

// KillContainer terminates and removes a containerID
// DELETE /v0/containers/kill/abc33412adqw
func KillContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// TODO: auth

	// Validate ContainerID
	containerID := params.ByName("id")
	if !vContainerID.MatchString(containerID) {
		response.WriteError(res, http.StatusBadRequest, er.InvalidContainerID)
		return
	}

	var userKey string
	cookie, err := req.Cookie("dock_session")
	if err != nil {
		fmt.Println("Error getting cookie")
		response.WriteError(res, http.StatusUnauthorized, "Not authorized")
		return
	}
	userKey = cookie.Value
	if userKey == "" {
		fmt.Println("cookie value is empty")
		response.WriteError(res, http.StatusUnauthorized, "Not authorized")
		return
	}
	userID := session.StripUserKey(userKey)
	var exists bool
	exists, err = session.ExistsRunConfig(userID)
	if err != nil || !exists {
		fmt.Println("session does not exist")
		response.WriteError(res, http.StatusUnauthorized, "Not authorized")
		return
	}
	var cfg dc.RunConfig
	cfg, err = session.GetRunConfig(userID)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}

	// XXX: Remove Container performs remove --force. Previous steps are not required.
	var port int
	port, err = strconv.Atoi(cfg.Port)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}

	err = dc.RemoveContainer(containerID, port)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}

	res.Write(response.Marshal())
}
