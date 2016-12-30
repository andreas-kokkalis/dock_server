package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
	"github.com/julienschmidt/httprouter"
)

// AdminKillContainer terminates and removes a containerID
// DELETE /v0/containers/kill/abc33412adqw
func AdminKillContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Validate ContainerID
	containerID := params.ByName("id")
	if !vContainerID.MatchString(containerID) {
		response.WriteError(res, http.StatusBadRequest, er.InvalidContainerID)
		return
	}
	// Get the cookie to get the admin key
	cookie, err := req.Cookie("ses")
	var cfg dc.RunConfig
	cfg, err = dc.GetAdminRunConfig(cookie.Value)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ContainerAlreadyKilled)
		return
	}
	// Kill containerID - 	// XXX: issues with deleting container
	port, _ := strconv.Atoi(cfg.Port)
	err = dc.RemoveContainer(containerID, port)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	// Remove Redis key
	err = dc.DeleteAdminRunConfig(cookie.Value)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	// response.SetData("Container deleted")
	fmt.Println("Returning successfully")
	response.WriteError(res, http.StatusOK, "There is no error")
	fmt.Println("I wrote the response error and I will try to return")
	return
}
