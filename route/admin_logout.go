package route

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/julienschmidt/httprouter"
)

// AdminLogout logs out an admin
func AdminLogout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	fmt.Println("request to logout")

	// Get session cookie
	cookie, err := req.Cookie("ses")
	if err != nil {
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Check if session exists in Redis. If it doesn't exist sent Unauthorized. Frontend will redirect to login page.
	err = dc.DeleteAdminSession(cookie.Value)
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
