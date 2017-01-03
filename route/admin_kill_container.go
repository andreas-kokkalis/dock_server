package route

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/route/er"
	"github.com/julienschmidt/httprouter"
)

// AdminKillContainer terminates and removes a containerID
// DELETE /v0/containers/kill/abc33412adqw
func AdminKillContainer(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	res.Header().Set("Content-Type", "application/json")
	response := NewResponse()

	// Validate ContainerID
	containerID := params.ByName("id")
	if !vContainerID.MatchString(containerID) {
		response.WriteError(res, http.StatusBadRequest, er.InvalidContainerID)
		return
	}
	// Get the cookie to get the admin key

	cookie, err := req.Cookie("ses")
	var cfg dc.RunConfig
	cfg, err = dc.GetAdminRunConfig(cookie.Value)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, er.ContainerAlreadyKilled)
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
	err = dc.RemoveContainer(containerID, port)
	if err != nil {
		fmt.Println(err.Error())
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Waited 100ms")

	// Remove Redis key
	err = dc.DeleteAdminRunConfig(cookie.Value)
	if err != nil {
		response.WriteError(res, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("%+v\n", req)
	fmt.Printf("\n\n%+v\n", res.Header())
	// defer res.WriteHeader(200)
	return
}
