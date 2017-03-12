package container

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	db     *store.DB
	redis  *store.RedisRepo
	docker *docker.Repo
	mapper *docker.PortMapper
}

// NewService creates a new Image Service
func NewService(db *store.DB, redis *store.RedisRepo, docker *docker.Repo, mapper *docker.PortMapper) Service {
	return Service{db, redis, docker, mapper}
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
func AdminRunContainer(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// Validate imageID
		imageID := params.ByName("id")
		if !api.VImageID.MatchString(imageID) {
			response.WriteError(res, http.StatusBadRequest, api.ErrInvalidImageID)
			return
		}

		cookie, _ := req.Cookie("ses")
		log.Println(cookie.Value)
		sessionExists, err := s.redis.ExistsAdminRunConfig(cookie.Value)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
			return
		}
		var cfg api.RunConfig
		username := "guest"
		password := "password"
		if sessionExists {
			cfg, err = s.redis.GetAdminRunConfig(cookie.Value)
			if err != nil {
				response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
				return
			}
			fmt.Printf("exists: %v\n", cfg)
			// Update the TTL
			err = s.redis.SetAdminRunConfig(cookie.Value, cfg)
			if err != nil {
				response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
				return
			}
		} else {
			// Session didn't exist. Admin requested to run a container for the first time.
			// Run container and set the session.
			var port int
			port, err = s.mapper.Reserve()
			if err != nil {
				log.Printf("[CreateContainer]: %v", err.Error())
				response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
				return
			}
			if port == -1 {
				log.Printf("[CreateContainer]: No ports were available to reserve.\n")
				response.WriteError(res, http.StatusInternalServerError, "there are no resources available in the system")
			}
			// Run the container and get the url
			cfg, err = s.docker.RunContainer(imageID, username, password, port)
			if err != nil {
				s.mapper.Remove(port)
				response.WriteError(res, http.StatusInternalServerError, err.Error())
				return
			}

			err = s.redis.SetAdminRunConfig(cookie.Value, cfg)
			if err != nil {
				// XXX: not sure if this is needed here, cause there was no error creating the cotnainer
				// s.mapper.Remove(port)
				response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
				return
			}
		}
		response.SetData(runResponse{URL: cfg.URL, Username: username, Password: password, ContainerID: cfg.ContainerID})
		_, _ = res.Write(response.Marshal())
	}
}

// AdminKillContainer terminates and removes a containerID
// DELETE /v0/containers/kill/abc33412adqw
func AdminKillContainer(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// Validate ContainerID
		containerID := params.ByName("id")
		if !api.VContainerID.MatchString(containerID) {
			response.WriteError(res, http.StatusBadRequest, api.ErrInvalidContainerID)
			return
		}
		// Get the cookie to get the admin key

		cookie, err := req.Cookie("ses")
		var cfg api.RunConfig
		cfg, err = s.redis.GetAdminRunConfig(cookie.Value)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, api.ErrContainerAlreadyKilled)
			return
		}
		fmt.Println(cfg)

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
			fmt.Println(err.Error())
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}
		time.Sleep(time.Millisecond * 100)
		fmt.Println("Waited 100ms")

		// Remove Redis key
		err = s.redis.DeleteAdminRunConfig(cookie.Value)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Printf("%+v\n", req)
		fmt.Printf("\n\n%+v\n", res.Header())
		// defer res.WriteHeader(200)
		return
	}
}

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
func CommitContainer(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// Parse containerID
		containerID := params.ByName("id")
		if !api.VContainerID.MatchString(containerID) {
			response.WriteError(res, http.StatusBadRequest, api.ErrInvalidContainerID)
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
			response.WriteError(res, http.StatusUnprocessableEntity, api.ErrInvalidPostData)
			return
		}
		// Create the new image
		err = s.docker.CommitContainer(data.Comment, data.Author, containerID, data.RefTag)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}
		// Get the cookie to get the admin key
		cookie, err := req.Cookie("ses")
		var cfg api.RunConfig
		cfg, err = s.redis.GetAdminRunConfig(cookie.Value)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
			return
		}
		if cfg.ContainerID == "" {
			response.WriteError(res, http.StatusBadRequest, api.ErrContainerAlreadyKilled)
			return
		}
		// Kill containerID - // XXX: issues with deleting container
		port, _ := strconv.Atoi(cfg.Port)
		// XXX: moved mapper.Remove here from RemoveContainer
		s.mapper.Remove(port)
		err = s.docker.ContainerRemove(containerID, port)
		if err != nil {
			fmt.Println(err.Error())
			response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
			return
		}
		// Remove Redis key
		err = s.redis.DeleteAdminRunConfig(cookie.Value)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, api.ErrServerError)
			return
		}
		log.Printf("[RT-CommitContainer]: attempting to write the response")
		response.SetData("A new image has been successfully created.")
		response.SetStatus(http.StatusOK, res)
		_, _ = res.Write(response.Marshal())
	}
}

// GetContainers returns list of containers by status.
// GET /v0/containers
// GET /v0/containers/:status
func GetContainers(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()

		// Validate if status has any of the accepted input
		status := params.ByName("status")
		if !api.VContainerState.MatchString(status) {
			// http.Error(res, er.InvalidContainerState, http.StatusUnprocessableEntity)
			response.AddError(api.ErrInvalidContainerState)
			response.SetStatus(http.StatusUnprocessableEntity, res)
			_, _ = res.Write(response.Marshal())
			return
		}

		// Get the list of containers
		containers, err := s.docker.GetContainers(status)
		if err != nil {
			// http.Error(res, er.ServerError, http.StatusInternalServerError)
			response.AddError(err.Error())
			response.SetStatus(http.StatusInternalServerError, res)
			_, _ = res.Write(response.Marshal())
			return
		}

		response.SetStatus(http.StatusOK, res)
		response.SetData(containers)
		_, _ = res.Write(response.Marshal())
	}
}
