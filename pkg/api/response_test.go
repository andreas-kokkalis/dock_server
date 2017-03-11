package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	assert := assert.New(t)
	res := NewResponse()
	assert.Equal(nil, res.Data, "It should be nil")
	assert.Equal("", res.Status, "It should be empty string")
	assert.Equal(0, len(res.Errors), "It should have 0 length")
}

func TestAddError(t *testing.T) {
	res := NewResponse()

	expected := "This is an error"
	res.AddError(expected)

	assert.Equal(t, expected, res.Errors[0], "It should be equal")
}

func TestSetData(t *testing.T) {
	res := NewResponse()

	expected := "test data"
	res.SetData(expected)

	assert.Equal(t, expected, res.Data, "It should be equal")
}

func TestMarshal(t *testing.T) {
	res := NewResponse()
	res.SetData("test data")

	expected := []byte(`{"data":"test data"}`)
	response := res.Marshal()

	assert.Equal(t, expected, response, "It should be equal")
}

/*
func TestSetStatus(t *testing.T) {
	res := NewResponse()
	var rw http.ResponseWriter

	res.SetStatus(http.StatusInternalServerError, rw)
	if res.Status != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("Expected: %v Got: %v", http.StatusInternalServerError, res.Status)
	}
}
*/
/*
func TestWriteError(t *testing.T) {
	res := NewResponse()
	var rw http.ResponseWriter
	res.WriteError(rw, http.StatusInternalServerError, errors.New("new error").Error())


}
*/
