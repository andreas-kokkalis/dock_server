package route

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
