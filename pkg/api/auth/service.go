package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/db"
	"github.com/jordic/lti"
	"github.com/julienschmidt/httprouter"
)

// Service for image
type Service struct {
	db    *db.DB
	redis *store.RedisRepo
}

// NewService creates a new Image Service
func NewService(db *db.DB, redis *store.RedisRepo) Service {
	return Service{db, redis}
}

var vAdminCookieVal = regexp.MustCompile(`^(adm:[a-f0-9]{32})$`)

// SessionAuth performs validation before invoking the route
func SessionAuth(s Service, handler httprouter.Handle) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

		// Get session cookie
		cookie, err := req.Cookie("ses")
		if err != nil {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// Validate Cookie value
		if !vAdminCookieVal.Match([]byte(cookie.Value)) {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			_, _ = res.Write([]byte("Not authorized"))
			return
		}

		// Check if session exists in Redis. If it doesn't exist sent Unauthorized. Frontend will redirect to login page.
		exists, err := s.redis.ExistsAdminSession(cookie.Value)
		if err != nil || !exists {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		handler(res, req, params)
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
	return func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		res.Header().Set("Content-Type", "application/json")
		response := api.NewResponse()
		log.Println("started login")
		// Parse post params
		decoder := json.NewDecoder(req.Body)
		var data loginRequest
		err := decoder.Decode(&data)
		if err != nil {
			response.WriteError(res, http.StatusUnprocessableEntity, err.Error())
			return
		}

		log.Println("decoded login request")
		log.Println(data.Username, data.Password)
		// Query the database and check if user exists
		var id int
		var password string
		row := s.db.QueryRow("SELECT id, password FROM admins WHERE username = $1", data.Username)
		err = row.Scan(&id, &password)
		switch {
		case err == sql.ErrNoRows:
			// Case when user does not exist in the database
			response.WriteError(res, http.StatusUnauthorized, api.ErrUsernameNotExists)
			return
		case err != nil:
			// Database error
			log.Println(err)
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}

		log.Println("got value from db")

		// Verify that passwords match
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(data.Password))
		if err != nil {
			response.WriteError(res, http.StatusUnauthorized, api.ErrPasswordMismatch)
			return
		}

		// Check whether the session exists or not.
		var sessionExists bool
		key := s.redis.CreateAdminKey(id)
		sessionExists, err = s.redis.ExistsAdminSession(key)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
			return
		}
		// Case session does not exist
		if !sessionExists {
			err = s.redis.SetAdminSession(key)
			if err != nil {
				response.WriteError(res, http.StatusInternalServerError, err.Error())
				return
			}
		}
		log.Println("creating cookie")
		// Whether the session exists or not, write the cookie
		cookie := &http.Cookie{
			Name:    "ses",
			Value:   key,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour),
		}
		http.SetCookie(res, cookie)
		fmt.Println(cookie)
		response.Data = "SUCCESS"
		_, _ = res.Write(response.Marshal())
		return
	}
}

// AdminLogout logs out an admin
func AdminLogout(s Service) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

		// fmt.Println("request to logout")

		// Get session cookie
		cookie, err := req.Cookie("ses")
		if err != nil {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Check if session exists in Redis. If it doesn't exist sent Unauthorized. Frontend will redirect to login page.
		err = s.redis.DeleteAdminSession(cookie.Value)
		if err != nil {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		cookie = &http.Cookie{
			Name:    "ses",
			Value:   "",
			Path:    "/",
			Expires: time.Now(),
		}
		http.SetCookie(res, cookie)
	}
}

const (
	oauthKey    = "oauth_key"
	oauthSecret = "oauth_secret"
)

// OAuth stuff
func OAuth(s Service, handler httprouter.Handle) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

		// TODO: Add check if user is not a student, return error

		// OAuth authentication of the TP requires to match the
		// request URL with the expected path. Since image IDs
		// change all the time, the path is constructed using
		// the imageID as extracted from the HTTP Header.
		path := fmt.Sprintf("https://%s%s", req.Host, req.URL.Path)
		fmt.Println(path)
		p := lti.NewProvider(oauthSecret, path)
		p.ConsumerKey = oauthKey

		ok, err := p.IsValid(req)
		if !ok {
			log.Println("invalid")
			fmt.Fprintf(res, "Invalid request...")
			return
		}
		if err != nil {
			log.Printf("Invalid request %s", err)
			return
		}
		handler(res, req, params)
	}
}
