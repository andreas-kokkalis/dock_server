package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	res := NewResponse()
	assert.Equal(nil, res.Data, "It should be nil")
	assert.Equal("", res.Status, "It should be empty string")
	assert.Equal(0, len(res.Errors), "It should have 0 length")
}

func TestAddError(t *testing.T) {
	t.Parallel()
	res := NewResponse()

	expected := "This is an error"
	res.AddError(expected)

	assert.Equal(t, expected, res.Errors[0], "It should be equal")
}

func TestSetData(t *testing.T) {
	t.Parallel()
	res := NewResponse()

	expected := "test data"
	res.SetData(expected)

	assert.Equal(t, expected, res.Data, "It should be equal")
}

func TestMarshal(t *testing.T) {
	t.Parallel()
	res := NewResponse()
	res.SetData("test data")

	expected := []byte(`{"data":"test data"}`)
	response := res.Marshal()

	assert.Equal(t, expected, response, "It should be equal")
}

func TestSetStatus(t *testing.T) {
	t.Parallel()
	res := NewResponse()
	handler := func(w http.ResponseWriter, r *http.Request) {
		res.SetStatus(http.StatusInternalServerError, w)
	}
	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestWriteError(t *testing.T) {
	t.Parallel()
	res := NewResponse()
	handler := func(w http.ResponseWriter, r *http.Request) {
		res.WriteError(w, http.StatusInternalServerError, "Error")
	}
	req := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Error", res.Errors[0])
}

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
