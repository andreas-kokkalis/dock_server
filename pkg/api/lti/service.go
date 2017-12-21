package lti

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/docker"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	db     *postgres.DB
	redis  *store.RedisRepo
	docker store.DockerRepository
	mapper *docker.PortMapper
}

// NewService creates a new Image Service
func NewService(db *postgres.DB, redis *store.RedisRepo, docker store.DockerRepository, mapper *docker.PortMapper) Service {
	return Service{db, redis, docker, mapper}
}

// Launch launches a url by imageID
// validate imageID
// extract user session
// check if container is running for that session
//	-- true: return current session
//  -- false: run container and return new session
func Launch(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

		fmt.Printf("Header: %+v\n", req.Header)
		fmt.Printf("Body:  %+v\n", req.Body)

		t, _ := template.ParseFiles("templates/html/assignment.html")
		// Validate imageID
		imageID := params.ByName("id")
		if !api.VImageID.MatchString(imageID) {
			_ = t.Execute(res, Resp{Error: "Invalid URL. Contact the administrator"})
		}

		// Parse LTI Post params
		err := req.ParseForm()
		if err != nil {
			_ = t.Execute(res, Resp{Error: "Invalid URL. Contact the administrator"})
		}
		// extract Canvas userID and store is as session key
		userID := req.PostFormValue("user_id")
		var sessionExists bool
		sessionExists, err = s.redis.ExistsUserRunConfig(userID)
		if err != nil {
			_ = t.Execute(res, Resp{Error: api.ErrServerError})
		}

		var cfg api.RunConfig
		if sessionExists {
			cfg, err = s.redis.GetUserRunConfig(userID)
			if err != nil {
				_ = t.Execute(res, Resp{Error: api.ErrServerError})
			}
			fmt.Printf("exists: %v\n", cfg)
			// Update the TTL
			err = s.redis.SetUserRunConfig(userID, cfg)
			if err != nil {
				_ = t.Execute(res, Resp{Error: api.ErrServerError})
			}
		} else {
			// SESSION didn'texist
			// Generate username and password
			username := "guest"
			// username := "canvas"
			password := "password"
			// password := newPassword()
			// Run container request
			var port int
			port, err = s.mapper.Reserve()
			if err != nil {
				log.Printf("[CreateContainer]: %v", err.Error())
				_ = t.Execute(res, Resp{Error: api.ErrServerError})
			}
			if port == -1 {
				log.Printf("[CreateContainer]: No ports were available to reserve.\n")
				_ = t.Execute(res, Resp{Error: "there are no resources available in the system"})
			}
			cfg, err = s.docker.ContainerRun(imageID, username, password, port)
			if err != nil {
				fmt.Println(err.Error())
				s.mapper.Remove(port)
				_ = t.Execute(res, Resp{Error: api.ErrServerError})
			}
			fmt.Printf("not exists: %v\n", cfg)
			// Set session
			err = s.redis.SetUserRunConfig(userID, cfg)
			if err != nil {
				// XXX: container is running
				// s.mapper.Remove(port)
				_ = t.Execute(res, Resp{Error: api.ErrServerError})
			}
		}

		// Whether the session exists or not, write the cookie
		cookie := &http.Cookie{
			Name:    "dock_session",
			Value:   s.redis.GetUserRunKey(userID),
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(res, cookie)
		fmt.Println(cookie)

		// Return HTML template with data
		_ = t.Execute(res, getResp(cfg))
	}
}

func getResp(cfg api.RunConfig) Resp {
	return Resp{
		ContainerID: cfg.ContainerID,
		Port:        cfg.Port,
		Username:    cfg.Username,
		Password:    cfg.Password,
		URL:         cfg.URL,
	}
}

// Resp ...
type Resp struct {
	ContainerID string `json:"id"`
	Port        string `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	URL         string `json:"url"`
	Error       string
}

/*
// KillContainer terminates and removes a containerID
// DELETE /v0/containers/kill/abc33412adqw
func KillContainer(s Service) httprouter.Handle {
    return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
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
	userID := dc.StripSessionKeyPrefix(cookieVal)
	var exists bool
	exists, err = dc.ExistsUserRunConfig(userID)
	if err != nil || !exists {
		fmt.Println("session does not exist")
		response.WriteError(res, http.StatusUnauthorized, "Not authorized")
		return
	}
	// Get session from Redis
	var cfg dc.RunConfig
	cfg, err = dc.GetUserRunConfig(userID)
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
	err = dc.DeleteUserRunConfig(userID)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteError(res, http.StatusInternalServerError, er.ServerError)
		return
	}

	res.Write(response.Marshal())
}
}
*/
