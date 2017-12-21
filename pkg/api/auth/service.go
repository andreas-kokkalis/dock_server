package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
	"github.com/jordic/lti"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	adminRepo repositories.AdminDBRepository
	redis     repositories.RedisRepository
}

// NewService creates a new Image Service
func NewService(adminRepo repositories.AdminDBRepository, redis repositories.RedisRepository) Service {
	return Service{adminRepo, redis}
}

var vAdminCookieVal = regexp.MustCompile(`^(adm:[a-f0-9]{32})$`)

// SessionAuth performs validation before invoking the route
// TODO: perhaps redirect to login service from here
func SessionAuth(s Service, handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

		// Get session cookie
		cookie, err := r.Cookie("ses")
		if err != nil {
			api.WriteErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		// Validate Cookie value
		if !vAdminCookieVal.Match([]byte(cookie.Value)) {
			api.WriteErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		// Check if session exists in Redis. If it doesn't exist sent Unauthorized. Frontend will redirect to login page.
		exists, err := s.redis.AdminSessionExists(cookie.Value)
		if err != nil || !exists {
			api.WriteErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		handler(w, r, params)
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminLogin authenticates an admin user
// Method POST
// Params: username, password
// TODO: validation of username and password
func AdminLogin(s Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Parse post params
		decoder := json.NewDecoder(r.Body)
		var data api.Admin
		err := decoder.Decode(&data)
		if err != nil {
			api.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Query the database and check if user exists
		admin, err := s.adminRepo.GetAdminByUsername(data)
		switch {
		case err == postgres.ErrNoResult:
			// Case when user does not exist in the database
			api.WriteErrorResponse(w, http.StatusUnauthorized, api.ErrUsernameNotExists)
			return
		case err != nil:
			// Database error
			api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Verify that passwords match
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(data.Password))
		if err != nil {
			api.WriteErrorResponse(w, http.StatusUnauthorized, api.ErrPasswordMismatch)
			return
		}

		// Check whether the session exists or not.
		var sessionExists bool
		key := s.redis.AdminSessionKeyCreate(admin.ID)
		if sessionExists, err = s.redis.AdminSessionExists(key); err != nil {
			api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Case session does not exist
		if !sessionExists {
			if err = s.redis.AdminSessionSet(key); err != nil {
				api.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		// Whether the session exists or not, write the cookie
		cookie := &http.Cookie{
			Name:    "ses",
			Value:   key,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(w, cookie)

		api.WriteOKResponse(w, "SUCCESS")
	}
}

// AdminLogout logs out an admin
func AdminLogout(s Service) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		// Get session cookie
		cookie, err := r.Cookie("ses")
		if err != nil {
			api.WriteErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		// Check if session exists in Redis. If it doesn't exist sent Unauthorized. Frontend will redirect to login page.
		err = s.redis.AdminSessionDelete(cookie.Value)
		if err != nil {
			api.WriteErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		cookie = &http.Cookie{
			Name:    "ses",
			Value:   "",
			Path:    "/",
			Expires: time.Now(),
		}
		http.SetCookie(w, cookie)
		api.WriteOKResponse(w, nil)
	}
}

const (
	oauthKey    = "oauth_key"
	oauthSecret = "oauth_secret"
)

// OAuth middleware
func OAuth(s Service, handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

		// OAuth authentication of the TP requires to match the
		// request URL with the expected path. Since image IDs
		// change all the time, the path is constructed using
		// the imageID as extracted from the HTTP Header.
		path := fmt.Sprintf("https://%s%s", r.Host, r.URL.Path)
		fmt.Println(path)
		p := lti.NewProvider(oauthSecret, path)
		p.ConsumerKey = oauthKey

		ok, err := p.IsValid(r)
		if !ok {
			log.Println("invalid")
			fmt.Fprintf(w, "Invalid request...")
			return
		}
		if err != nil {
			log.Printf("Invalid request %s", err)
			return
		}
		handler(w, r, params)
	}
}
