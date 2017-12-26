package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	type Data struct {
		Foo string `json:"foo"`
	}
	bodyData := Data{Foo: "bar"}
	// Since the request can only be used within a gomega assertion,
	// register the fail handler, and a recover function
	gomega.RegisterFailHandler(ginkgo.Fail)
	defer ginkgo.GinkgoRecover()

	response := NewResponse(200, `{"foo":"bar"}`)
	response.recorder = httptest.NewRecorder()
	api.WriteOKResponse(response.recorder, bodyData)

	var actualData Data
	response.Unmarshall(&actualData)
	assert.Equal(t, response.recorder.Code, response.expectedCode)
	assert.Equal(t, bodyData, actualData)
}
