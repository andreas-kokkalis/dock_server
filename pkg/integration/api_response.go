package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/onsi/gomega"
)

// Response struct for asserting an http API response
type Response struct {
	expectedCode int
	expectedBody string
	recorder     *httptest.ResponseRecorder
}

// NewResponse initializes a Response object hat is used to test the expected output against the actual HTTP response
func NewResponse(expectedCode int, expectedJSONBody string) *Response {
	return &Response{
		expectedCode: expectedCode,
		expectedBody: expectedJSONBody,
		recorder:     httptest.NewRecorder(),
	}
}

// ToString returns the JSON API response
func (r *Response) ToString() string {
	return r.recorder.Body.String()
}

// Unmarshall takes a JSON api response and unmarshals the result data into the target interface
func (r *Response) Unmarshall(target interface{}) {
	var res api.Response
	err := json.Unmarshal(r.recorder.Body.Bytes(), &res)
	gomega.Expect(err).To(gomega.BeNil(), "error unmarshalling JSON response")
	byteData, _ := json.Marshal(res.Data)
	err = json.Unmarshal(byteData, target)
	gomega.Expect(err).To(gomega.BeNil(), "error unmarshalling JSON data to target datastructure")
}

// nolint
type Res struct {
	Code    int          `json:"Code"`
	Headers http.Header  `json:"Headers"`
	Body    api.Response `json:"Body,omitempty"`
}

// nolint
type LogResponse struct {
	HTTPResponse Res `json:"HTTP_Response"`
}

// pretty function of Respone logs the API response for debugging purposes
func (r *Response) pretty() string {

	var body api.Response
	_ = json.Unmarshal(r.recorder.Body.Bytes(), &body)

	lr := LogResponse{Res{
		Code:    r.recorder.Code,
		Headers: r.recorder.Header(),
		Body:    body,
	}}
	res, err := json.MarshalIndent(lr, "", "  ")
	gomega.Expect(err).To(gomega.BeNil())
	return string(res)
}

// Code returns the actual HTTP status Code
func (r *Response) Code() int {
	return r.recorder.Code
}
