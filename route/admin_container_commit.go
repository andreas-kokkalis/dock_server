package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
	"github.com/julienschmidt/httprouter"
)

type commitContainerRequest struct {
	Comment string `json:"comment"`
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

	// Parse containerID
	containerID := params.ByName("id")
	if !vContainerID.MatchString(containerID) {
		response.WriteError(res, http.StatusBadRequest, er.InvalidContainerID)
		return
	}
	// Parse post params
	decoder := json.NewDecoder(req.Body)
	var data commitContainerRequest
	err := decoder.Decode(&data)
	if err != nil {
		response.WriteError(res, http.StatusUnprocessableEntity, err.Error())
		return
	}
	// Validate post params
	if data.Comment == "" || data.Author == "" || data.RefTag == "" {
		response.WriteError(res, http.StatusUnprocessableEntity, er.InvalidPostData)
		return
	}
	// Create the new image
	err = dc.CommitContainer(data.Comment, data.Author, containerID, data.RefTag)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	// Get the cookie to get the admin key
	cookie, err := req.Cookie("ses")
	var cfg dc.RunConfig
	cfg, err = dc.GetAdminRunConfig(cookie.Value)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}
	if cfg.ContainerID == "" {
		response.WriteError(res, http.StatusBadRequest, er.ContainerAlreadyKilled)
		return
	}
	// Kill containerID - // XXX: issues with deleting container
	port, _ := strconv.Atoi(cfg.Port)
	err = dc.RemoveContainer(containerID, port)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}
	// Remove Redis key
	err = dc.DeleteAdminRunConfig(cookie.Value)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}
	log.Printf("[RT-CommitContainer]: attempting to write the response")
	response.SetData("A new image has been successfully created.")
	response.SetStatus(http.StatusOK, res)
	res.Write(response.Marshal())
}
