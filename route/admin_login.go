package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/andreas-kokkalis/dock-server/db"
	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
	"github.com/julienschmidt/httprouter"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AdminLogin authenticates an admin user
// Method POST
// Params: username, password
// TODO: validation of username and password
func AdminLogin(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Parse post params
	decoder := json.NewDecoder(req.Body)
	var data loginRequest
	err := decoder.Decode(&data)
	if err != nil {
		response.WriteError(res, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Query the database and check if user exists
	var id int
	var password string
	row := db.PG.QueryRow("SELECT id, password FROM admins WHERE username = $1", data.Username)
	err = row.Scan(&id, &password)
	switch {
	case err == sql.ErrNoRows:
		// Case when user does not exist in the database
		response.WriteError(res, http.StatusUnauthorized, er.UsernameNotExists)
		return
	case err != nil:
		// Database error
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}

	// Verify that passwords match
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(data.Password))
	if err != nil {
		response.WriteError(res, http.StatusUnauthorized, er.PasswordMismatch)
		return
	}

	// Check whether the session exists or not.
	var sessionExists bool
	key := dc.CreateAdminKey(id)
	sessionExists, err = dc.ExistsAdminSession(key)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	// Case session does not exist
	if !sessionExists {
		err = dc.SetAdminSession(key)
		if err != nil {
			response.WriteError(res, http.StatusInternalServerError, err.Error())
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
	http.SetCookie(res, cookie)
	fmt.Println(cookie)
	response.Data = "SUCCESS"
	res.Write(response.Marshal())
	return
}
