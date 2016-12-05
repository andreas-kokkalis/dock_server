package route

import (
	"net/http"

	"github.com/andreas-kokkalis/dock-server/session"
	"github.com/julienschmidt/httprouter"
)

// AuthAdmin performs validation before invoking the route
func AuthAdmin(handler httprouter.Handle) httprouter.Handle {
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
			res.Write([]byte("Not authorized"))
			return
		}

		// Check if session exists in Redis. If it doesn't exist sent Unauthorized. Frontend will redirect to login page.
		exists := false
		exists, err = session.ExistsAdminSession(cookie.Value)
		if err != nil || !exists {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		handler(res, req, params)
	}
}
