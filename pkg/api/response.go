package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response of the API
type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
	Status string      `json:"status,omitempty"`
}

// NewResponse returns a new Response struct
func NewResponse() *Response {
	return &Response{}
}

// AddError adds an error message in the slice
func (r *Response) AddError(err string) {
	r.Errors = append(r.Errors, err)
}

// WriteError writes in http Response the statusCode and the error message in response
func (r *Response) WriteError(res http.ResponseWriter, statusCode int, err string) {
	res.WriteHeader(statusCode)
	r.AddError(err)
	_, _ = res.Write(r.Marshal())
	log.Println(err)
}

// SetStatus will set the http status in the response
func (r *Response) SetStatus(statusCode int, res http.ResponseWriter) {
	r.Status = http.StatusText(statusCode)
	res.WriteHeader(statusCode)
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

// WriteErrorResponse writes an error
func WriteErrorResponse(w http.ResponseWriter, statusCode int, msg ...string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).
		Encode(&Response{
			Status: http.StatusText(statusCode),
			Errors: msg,
		})
}

// WriteOKResponse writes a valid response
func WriteOKResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&Response{Data: data})
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
