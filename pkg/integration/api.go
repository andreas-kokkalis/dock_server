package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/onsi/gomega"
)

// Request struct for performing an http request
type Request struct {
	Method string
	Target string
	Body   interface{}
}

// Log function of Request logs the request for debugging purposes
func (r Request) Log() string {
	var body []byte
	if r.Body != nil {
		body, _ = json.MarshalIndent(r.Body, "", "  ")
		return fmt.Sprintf(
			"\n-------------\n"+
				"HTTP Request\n"+
				"-------------\n"+
				"%s %s\n"+
				"%s\n\n",
			r.Method, r.Target,
			string(body))
	}
	return fmt.Sprintf(
		"\n------------\n"+
			"HTTP Request\n"+
			"------------\n"+
			"%s %s\n\n",
		r.Method, r.Target)
}

// Response struct for asserting an http API response
type Response struct {
	ExpectedCode   int
	ExpectedBody   string
	actualResponse *httptest.ResponseRecorder
}

// ToString returns the JSON API response
func (r *Response) ToString() string {
	return r.actualResponse.Body.String()
}

// ToStructure unmarshals the JSON API response to the target data structure
// target should be a pointer to a structure
func (r *Response) ToStructure(target interface{}) {
	err := json.Unmarshal(r.actualResponse.Body.Bytes(), target)
	gomega.Expect(err).To(gomega.BeNil(), "error unmarshalling JSON response to data structure")
}

// Log function of Respone logs the API response for debugging purposes
func (r Response) Log() string {
	var apiOut bytes.Buffer
	_ = json.Indent(&apiOut, r.actualResponse.Body.Bytes(), "", "  ")
	return fmt.Sprintf(
		"\n-----------\n"+
			"API Response\n"+
			"-----------\n"+
			"Status Code: %d\n"+
			"%s\n\n",
		r.ExpectedCode, apiOut.String())
}

// Code returns the actual HTTP status Code
func (r Response) Code() int {
	return r.actualResponse.Code
}

// performRequest performs an HTTP request to an API endpoint and records the response to
// an httptest.ResponseRecorder. The recorder is returned to be used in assertions.
func performRequest(router http.Handler, request Request, response *Response) {
	var reqBody io.ReadCloser

	if request.Body != nil {
		jsonBody, err := json.Marshal(request.Body)
		gomega.Expect(err).To(gomega.BeNil())
		reqBody = ioutil.NopCloser(bytes.NewReader(jsonBody))
	} else {
		reqBody = nil
	}

	req := httptest.NewRequest(request.Method, request.Target, reqBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	response.actualResponse = w
}

// EvalAPIResponse performs an HTTP request, records the output and asserts if it matches against the expected response code and body.
func EvalAPIResponse(router http.Handler, request Request, response *Response, topDir string, logger *log.Logger) {
	performRequest(router, request, response)
	logger.Println(request.Log())
	logger.Println(response.Log())
	gomega.Expect(response.Code()).To(gomega.Equal(response.ExpectedCode), "status codes do not match")

	diff, err := CompareRegexJSON(response.ExpectedBody, response.ToString(), topDir)
	gomega.Expect(err).To(gomega.BeNil(), "Diff tool returned error")
	gomega.Expect(diff).To(gomega.Equal(""), "Diff is not empty")
}
