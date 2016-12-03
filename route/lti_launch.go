package route

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/andreas-kokkalis/dock-server/dc"
	"github.com/andreas-kokkalis/dock-server/er"
	"github.com/andreas-kokkalis/dock-server/session"
	"github.com/jordic/lti"
	"github.com/julienschmidt/httprouter"
)

const (
	oauthKey    = "oauth_key"
	oauthSecret = "oauth_secret"
)

// OAuth stuff
func OAuth(handler httprouter.Handle) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

		/*
			err := req.ParseForm()
			if err != nil { // Handle error here via logging and then return }
			key := req.PostFormValue("oauth_consumer_key")
			secret := req.PostFormValue("oauth_signature")
			signatureMethod := req.PostFormValue("oauth_signature_method")
			user := req.PostFormValue("user_id")
			fmt.Printf("\n%s \n%s \n%s \n%s", key, secret, signatureMethod, user)
		*/

		// Provider requires to match the request url with the secret.
		// Since the request URL depends on imageID it should constuct it from the header
		path := fmt.Sprintf("https://%s%s", req.Host, req.URL.Path)
		fmt.Println(path)

		p := lti.NewProvider(oauthSecret, path)
		p.ConsumerKey = oauthKey

		ok, err := p.IsValid(req)
		if !ok {
			fmt.Fprintf(res, "Invalid request...")
		}
		if err != nil {
			log.Printf("Invalid request %s", err)
			return
		}
		handler(res, req, params)
	}
}

// LTILaunch launches a url by imageID
// validate imageID
// extract user session
// check if container is running for that session
//	-- true: return current session
//  -- false: run container and return new session
func LTILaunch(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	t, _ := template.ParseFiles("templates/html/assignment.html")
	// Validate imageID
	imageID := params.ByName("id")
	if !vImageID.MatchString(imageID) {
		t.Execute(res, Resp{Error: "Invalid URL. Contact the administrator"})
	}

	// Parse LTI Post params
	err := req.ParseForm()
	if err != nil {
		t.Execute(res, Resp{Error: "Invalid URL. Contact the administrator"})
	}
	// extract Canvas userID and store is as session key
	userID := req.PostFormValue("user_id")
	var sessionExists bool
	sessionExists, err = session.ExistsRunConfig(userID)
	if err != nil {
		t.Execute(res, Resp{Error: "Hello"})
	}

	var cfg dc.RunConfig
	if sessionExists {
		cfg, err = session.GetRunConfig(userID)
		if err != nil {
			// TODO: handle error
		}
		fmt.Printf("exists: %v\n", cfg)
		// Update the TTL
		err = session.SetRunConfig(userID, cfg)
		if err != nil {
			t.Execute(res, Resp{Error: er.ServerError})
		}
	} else {
		// SESSION didn'texist
		// Generate username and password
		username := newPassword()
		password := newPassword()
		// Run container request
		cfg, err = dc.RunContainer(imageID, username, password)
		if err != nil {
			fmt.Println(err.Error())
			t.Execute(res, Resp{Error: er.ServerError})
		}
		fmt.Printf("not exists: %v\n", cfg)
		// Set session
		err = session.SetRunConfig(userID, cfg)
		if err != nil {
			t.Execute(res, Resp{Error: er.ServerError})
		}
	}
	// Return HTML template with data
	t.Execute(res, Resp{Config: cfg})
}

// Resp ...
type Resp struct {
	Config dc.RunConfig
	Error  string
}

func newPassword() string {
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	length := 6
	newPword := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
			panic(err)
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}
			newPword[i] = chars[c%clen]
			i++
			if i == length {
				return string(newPword)
			}
		}
	}
}
