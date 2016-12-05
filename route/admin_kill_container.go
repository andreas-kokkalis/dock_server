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
	// Get session cookie
	var cookieVal string
	cookie, err := req.Cookie("dock_session")
	if err != nil {
		fmt.Println("Error getting cookie")
		response.WriteError(res, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	// Get cookie value
	cookieVal = cookie.Value
	if cookieVal == "" {
		fmt.Println("cookie value is empty")
		response.WriteError(res, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	// Check if session exists in Redis
	userID := session.StripKey(cookieVal)
	var exists bool
	exists, err = session.ExistsRunConfig(userID)
	if err != nil || !exists {
		fmt.Println("session does not exist")
		response.WriteError(res, http.StatusUnauthorized, "Not authorized")
		return
	}
	// Get session from Redis
	var cfg dc.RunConfig
	cfg, err = session.GetRunConfig(userID)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}
	// Prepare the port to user in Remove container call
	var port int
	port, err = strconv.Atoi(cfg.Port)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}
	// XXX: issues with deleting container
	err = dc.RemoveContainer(containerID, port)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}

	// Delete the user session
	err = session.DeleteRunConfig(userID)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}

	res.Write(response.Marshal())
}
