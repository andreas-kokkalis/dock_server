package api

import (
	"encoding/json"
	"net/http"
)

// Response of the API
type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
	Status string      `json:"status,omitempty"`
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
	_ = json.NewEncoder(w).Encode(&Response{Data: data})
}
