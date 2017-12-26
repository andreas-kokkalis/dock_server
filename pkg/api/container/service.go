package container

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/portmapper"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	redis  repositories.RedisRepository
	docker repositories.DockerRepository
	mapper *portmapper.PortMapper
}

// NewService creates a new Image Service
func NewService(redis repositories.RedisRepository, docker repositories.DockerRepository, mapper *portmapper.PortMapper) Service {
	return Service{redis, docker, mapper}
}

/*
type runRequest struct {
	Username string `json:"user"`
	Password string `json:"pwd"`
}
*/

type runResponse struct {
	URL         string `json:"url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ContainerID string `json:"id"`
}

// AdminRunContainer POST
// POST /v0/containers/run
func (s Service) AdminRunContainer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Validate imageID
	imageID := params.ByName("id")
	if !api.VImageID.MatchString(imageID) {
		api.WriteErrorResponse(w, http.StatusBadRequest, api.ErrInvalidImageID)
		return
	}

	cookie, err := r.Cookie("ses")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "admin session cookie is not set")
		return
	}
	sessionExists, err := s.redis.AdminRunConfigExists(cookie.Value)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	var cfg api.RunConfig
	username := "guest"
	password := "password"
	if sessionExists {
		cfg, err = s.redis.AdminRunConfigGet(cookie.Value)
		if err != nil {
			api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Printf("exists: %v\n", cfg)
		// Update the TTL
		err = s.redis.AdminRunConfigSet(cookie.Value, cfg)
		if err != nil {
			api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// Session didn't exist. Admin requested to run a container for the first time.
		// Run container and set the session.
		var port int
		port, err = s.mapper.Reserve()
		if err != nil {
			api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if port == -1 {
			log.Printf("[CreateContainer]: No ports were available to reserve.\n")
			api.WriteErrorResponse(w, http.StatusInternalServerError, "there are no resources available in the system")
			return
		}
		// Run the container and get the url
		cfg, err = s.docker.ContainerRun(imageID, username, password, port)
		if err != nil {
			s.mapper.Remove(port)
			api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = s.redis.AdminRunConfigSet(cookie.Value, cfg)
		if err != nil {
			// XXX: not sure if this is needed here, cause there was no error creating the cotnainer
			// s.mapper.Remove(port)
			api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	api.WriteOKResponse(w, api.ContainerRun{URL: cfg.URL, Username: username, Password: password, ContainerID: cfg.ContainerID})
}

// AdminKillContainer terminates and removes a containerID
// DELETE /v0/containers/kill/abc33412adqw
func (s Service) AdminKillContainer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Validate ContainerID
	containerID := params.ByName("id")
	if !api.VContainerID.MatchString(containerID) {
		api.WriteErrorResponse(w, http.StatusBadRequest, api.ErrInvalidContainerID)
		return
	}
	// Get the cookie to get the admin key
	cookie, err := r.Cookie("ses")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "admin session cookie is not set")
		return
	}
	var cfg api.RunConfig
	cfg, err = s.redis.AdminRunConfigGet(cookie.Value)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, api.ErrContainerAlreadyKilled)
		return
	}
	// fmt.Println(cfg)

	// XXX: There is an issue when terminating the TLS handshake in the running container.
	// The frontend client produces a network_changed error.

	// Kill containerID - 	// XXX: issues with deleting container
	/*
			XXX: This works

		curl -k --verbose --cookie "ses=adm:7ff10abb653dead4186089acbd2b7891" -X DELETE -H "Cache-Control: no-cache" "https:ll/ddaab041abc5"ntainers/kil
		*   Trying 127.0.0.1...
		* Connected to localhost (127.0.0.1) port 8080 (#0)
		* found 173 certificates in /etc/ssl/certs/ca-certificates.crt
		* found 697 certificates in /etc/ssl/certs
		* ALPN, offering http/1.1
		* SSL connection using TLS1.2 / ECDHE_ECDSA_AES_128_GCM_SHA256
		* 	 server certificate verification SKIPPED
		* 	 server certificate status verification SKIPPED
		* 	 common name: KTH (does not match 'localhost')
		* 	 server certificate expiration date OK
		* 	 server certificate activation date OK
		* 	 certificate public key: EC
		* 	 certificate version: #3
		* 	 subject: C=SE,ST=Sweden,L=Stockholm,O=KTH,OU=KTH,CN=KTH,EMAIL=andreas@kth.se
		* 	 start date: Sat, 26 Nov 2016 15:37:07 GMT
		* 	 expire date: Tue, 24 Nov 2026 15:37:07 GMT
		* 	 issuer: C=SE,ST=Sweden,L=Stockholm,O=KTH,OU=KTH,CN=KTH,EMAIL=andreas@kth.se
		* 	 compression: NULL
		* ALPN, server accepted to use http/1.1
		> DELETE /v0/admin/containers/kill/ddaab041abc5 HTTP/1.1
		> Host: localhost:8080
		> User-Agent: curl/7.47.0
		> Accept: *\/*
		> Cookie: ses=adm:7ff10abb653dead4186089acbd2b7891
		> Cache-Control: no-cache
		>
		< HTTP/1.1 200 OK
		< Content-Type: application/json
		< Date: Sun, 01 Jan 2017 18:44:30 GMT
		< Content-Length: 0
		<
		* Connection #0 to host localhost left intact
	*/

	port, _ := strconv.Atoi(cfg.Port)
	// XXX: moved remove from repo to here
	s.mapper.Remove(port)
	err = s.docker.ContainerRemove(containerID, port)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// time.Sleep(time.Millisecond * 100)
	// fmt.Println("Waited 100ms")

	// Remove Redis key
	err = s.redis.AdminRunConfigDelete(cookie.Value)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fmt.Printf("%+v\n", req)
	// fmt.Printf("\n\n%+v\n", res.Header())
	// defer res.WriteHeader(200)
	api.WriteOKResponse(w, "Container Killed")
}

// CommitContainer creates a new image out of a running container
// POST /v0/containers/commit/:id
// JSON data:
//	* Comment
//	* Author
//	* RefTag
func (s Service) CommitContainer(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Parse containerID
	containerID := params.ByName("id")
	if !api.VContainerID.MatchString(containerID) {
		api.WriteErrorResponse(w, http.StatusBadRequest, api.ErrInvalidContainerID)
		return
	}
	// Parse post params
	decoder := json.NewDecoder(r.Body)
	var data api.ContainerCommitRequest
	err := decoder.Decode(&data)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	// Validate post params
	if data.Comment == "" || data.Author == "" || data.RefTag == "" {
		api.WriteErrorResponse(w, http.StatusBadRequest, api.ErrInvalidPostData)
		return
	}
	// Create the new image
	newImgID, err := s.docker.ContainerCommit(data.Comment, data.Author, containerID, data.RefTag)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Get the cookie to get the admin key
	cookie, err := r.Cookie("ses")
	if err != nil {
		api.WriteErrorResponse(w, http.StatusBadRequest, "admin session cookie is not set")
		return
	}
	var cfg api.RunConfig
	cfg, err = s.redis.AdminRunConfigGet(cookie.Value)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if cfg.ContainerID == "" {
		api.WriteErrorResponse(w, http.StatusInternalServerError, api.ErrContainerAlreadyKilled)
		return
	}
	// Kill containerID - // XXX: issues with deleting container
	port, _ := strconv.Atoi(cfg.Port)
	// XXX: moved mapper.Remove here from RemoveContainer
	s.mapper.Remove(port)
	err = s.docker.ContainerRemove(containerID, port)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Remove Redis key
	err = s.redis.AdminRunConfigDelete(cookie.Value)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.WriteOKResponse(w, api.ContainerCommitResponse{ImageID: newImgID})
	log.Printf("[RT-CommitContainer]: attempting to write the response")
}

// GetContainers returns list of containers by status.
// GET /v0/containers
// GET /v0/containers/:status
func (s Service) GetContainers(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// Validate if status has any of the accepted input
	status := params.ByName("status")
	if !api.VContainerState.MatchString(status) {
		api.WriteErrorResponse(w, http.StatusBadRequest, api.ErrInvalidContainerState)
		return
	}

	// Get the list of containers
	containers, err := s.docker.ContainerList(status)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.WriteOKResponse(w, containers)
}
