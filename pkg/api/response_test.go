package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	WriteErrorResponse(w, http.StatusInternalServerError, "Error")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var res Response
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, "Error", res.Errors[0])
}

func TestWriteOKResponse(t *testing.T) {
	w := httptest.NewRecorder()
	WriteOKResponse(w, "Data")
	assert.Equal(t, http.StatusOK, w.Code)
	var res Response
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, "Data", res.Data.(string))
}
