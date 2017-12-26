package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/onsi/gomega"
)

// Request struct for performing an HTTP request
type Request struct {
	method string
	url    string
	body   interface{}
	// HTTPRequest models the HTTP request. It's Exposed in order to set custon request headers, cookies and authentication for a request.
	HTTPRequest *http.Request
}

// NewRequest intializes a Request object
func NewRequest(method, url string, body interface{}) *Request {
	jsonBody, err := json.Marshal(body)
	gomega.Expect(err).To(gomega.BeNil(), "Error marshaling body parameter to json")
	return &Request{
		method:      method,
		url:         url,
		body:        body,
		HTTPRequest: httptest.NewRequest(method, url, ioutil.NopCloser(bytes.NewReader(jsonBody))),
	}
}

// WithSessionCookie sets a session cookie in the HTTP request
func (r *Request) WithSessionCookie(val string) *Request {
	r.HTTPRequest.AddCookie(&http.Cookie{Name: "ses", Value: val})
	return r
}

// pretty pretty creates a pretty string that models an HTTP request
func (r *Request) pretty() string {
	lg := LogRequest{HTTPRequest: Req{r.method, r.url, r.body}}
	request, err := json.MarshalIndent(lg, "", "  ")
	gomega.Expect(err).To(gomega.BeNil(), "Error marshaling request to JSON")
	return string(request)
}

// nolint
type Req struct {
	Method string      `json:"Method"`
	URL    string      `json:"URL"`
	Body   interface{} `json:"Body,omitempty"`
}

// nolint
type LogRequest struct {
	HTTPRequest Req `json:"HTTP_Request"`
}
