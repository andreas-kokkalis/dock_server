package route

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Ping tralalo
func Ping(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Write([]byte("Pong"))
}
