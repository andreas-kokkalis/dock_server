package route

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/andreas-kokkalis/dock-server/session"
	"github.com/andreas-kokkalis/dock-server/srv"
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
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Query the database and check if user exists
	var id int
	var password string
	row := srv.PG.QueryRow("SELECT id, password FROM admins WHERE username = $1", data.Username)
	err = row.Scan(&id, &password)
	switch {
	case err == sql.ErrNoRows:
		// Case when user does not exist in the database
		response.AddError(er.UsernameNotExists)
		res.Write(response.Marshal())
		return
	case err != nil:
		// Database error
		http.Error(res, err.Error(), http.StatusInternalServerError)
		response.AddError(err.Error())
		res.Write(response.Marshal())
		return
	}

	// Verify that passwords match
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(data.Password))
	if err != nil {
		response.AddError(er.PasswordMismatch)
		res.Write(response.Marshal())
		return
	}

	// Case when user exists. Check if there is an existing session
	var sessionExists bool
	sessionExists, err = session.AdminExists(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		response.AddError(err.Error())
		res.Write(response.Marshal())
		return
	}
	// Case session does not exist
	if !sessionExists {
		err = session.AdminAdd(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			response.AddError(err.Error())
			res.Write(response.Marshal())
			return
		}
	}
	// Whether the session exists or not, write the cookie
	cookie := &http.Cookie{
		Name:    "ses",
		Value:   session.GetAdminKey(id),
		Domain:  "KTH",
		Expires: time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(res, cookie)
	fmt.Println(cookie)

	response.Data = "SUCCESS"
	res.Write(response.Marshal())
	return
}
