package route

import (
	"encoding/json"
	"net/http"
)

// Response of the API
type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
	Status string      `json:"status"`
}

// NewResponse returns a new Response struct
func NewResponse() *Response {
	return &Response{}
}

// AddError adds an error message in the slice
func (r *Response) AddError(err string) {
	r.Errors = append(r.Errors, err)
}

// SetStatus will set the http status in the response
func (r *Response) SetStatus(statusCode int) {
	r.Status = http.StatusText(statusCode)
}

// SetData will set the http status in the response
func (r *Response) SetData(data interface{}) {
	r.Data = data
}

// Marshal will perform JSON Marshal on the response and return a byte slice
// XXX: Ignores errors
func (r *Response) Marshal() []byte {
	js, _ := json.Marshal(r)
	return js
}
