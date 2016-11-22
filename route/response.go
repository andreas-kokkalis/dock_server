package route

import "encoding/json"

// Response are formatted like this
type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

// NewResponse returns a new Response struct
func NewResponse() *Response {
	return &Response{}
}

// AddError adds an error message in the slice
func (r *Response) AddError(err string) {
	r.Errors = append(r.Errors, err)
}

// Marshal will perform JSON Marshal on the response and return a byte slice
// XXX: Ignores errors
func (r *Response) Marshal() []byte {
	js, _ := json.Marshal(r)
	return js
}
