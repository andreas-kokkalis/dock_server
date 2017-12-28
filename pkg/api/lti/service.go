package lti

import (
	"html/template"
	"net/http"
	"path"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/portmapper"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	redis       repositories.RedisRepository
	docker      repositories.DockerRepository
	mapper      *portmapper.PortMapper
	templateDir string
}

// NewService creates a new Image Service
func NewService(redis repositories.RedisRepository, docker repositories.DockerRepository, mapper *portmapper.PortMapper, templateDir string) Service {
	return Service{redis, docker, mapper, templateDir}
}

// nolint
const (
	ErrInvalidImageID        = "Invalid ImageID. The Launch URL has not been configured correctly"
	ErrInvalidFormData       = "Invalid form Data. There is an issue with the LTI integration"
	ErrResourceQuotaExceeded = "there are no resources available in the system"
	ErrContainerRun          = "unable to run container"
)

// Launch launches a url by imageID
// validate imageID
// extract user session
// check if container is running for that session
//	-- true: return current session
//  -- false: run container and return new session
func (s Service) Launch(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// fmt.Printf("Header: %+v\n", req.Header)
	// fmt.Printf("Body:  %+v\n", req.Body)
	tmp := path.Join(s.templateDir, "templates/html/assignment.html")
	t, err := template.ParseFiles(tmp)
	if err != nil {
		api.WriteErrorResponse(w, http.StatusInternalServerError, "Unable to find template for assignment")
		return
	}

	// Validate imageID
	imageID := params.ByName("id")
	if !api.VImageID.MatchString(imageID) {
		w.WriteHeader(http.StatusBadRequest)
		_ = t.Execute(w, Resp{Error: ErrInvalidImageID})
		return
	}

	// Parse LTI Post params
	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = t.Execute(w, Resp{Error: ErrInvalidFormData})
		return
	}

	// extract Canvas userID and store is as session key
	userID := r.PostFormValue("user_id")
	var sessionExists bool
	sessionExists, err = s.redis.UserRunConfigExists(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = t.Execute(w, Resp{Error: api.ErrServerError})
		return
	}

	var cfg api.RunConfig
	if sessionExists {
		cfg, err = s.redis.UserRunConfigGet(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = t.Execute(w, Resp{Error: api.ErrServerError})
			return
		}
		// fmt.Printf("exists: %v\n", cfg)
		// Update the TTL
		err = s.redis.UserRunConfigSet(userID, cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = t.Execute(w, Resp{Error: api.ErrServerError})
			return
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
		if err != nil || port == -1 {
			// log.Printf("[CreateContainer]: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_ = t.Execute(w, Resp{Error: ErrResourceQuotaExceeded})
			return
		}
		cfg, err = s.docker.ContainerRun(imageID, username, password, port)
		if err != nil {
			// fmt.Println(err.Error())
			s.mapper.Remove(port)
			w.WriteHeader(http.StatusInternalServerError)
			_ = t.Execute(w, Resp{Error: ErrContainerRun})
			return
		}
		// Set session
		err = s.redis.UserRunConfigSet(userID, cfg)
		if err != nil {
			// if key failed to set, contaienr will be killed by the mapper eventually
			w.WriteHeader(http.StatusInternalServerError)
			_ = t.Execute(w, Resp{Error: api.ErrServerError})
			return
		}
	}

	// Whether the session exists or not, write the cookie
	cookie := &http.Cookie{
		Name:    "dock_session",
		Value:   s.redis.UserRunKeyGet(userID),
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	// Return HTML template with data
	w.WriteHeader(http.StatusOK)
	_ = t.Execute(w, getResp(cfg))
	return
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
TODO

Terminate button for student doesn't work.
Add a call from the template to the terminate endpoint for student.
Container ID can be found from the session key with prefix dock_session


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
