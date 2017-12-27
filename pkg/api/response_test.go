package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteErrorResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		WriteErrorResponse(w, http.StatusInternalServerError, "Error")
	}
	w := httptest.NewRecorder()
	handler(w, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var res Response
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "Error", res.Errors[0])
}

func TestWriteOKResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		WriteOKResponse(w, "Data")
	}
	w := httptest.NewRecorder()
	handler(w, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusOK, w.Code)

	var res Response
	_ = json.Unmarshal(w.Body.Bytes(), &res)
	assert.Equal(t, "Data", res.Data.(string))
}
