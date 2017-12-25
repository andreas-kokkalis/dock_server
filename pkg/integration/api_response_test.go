package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	type Data struct {
		Foo string `json:"foo"`
	}

	d1 := Data{"bar"}

	tests := []struct {
		code int
		body string
		name string
	}{
		{
			http.StatusOK,
			`{"foo": "bar"}`,
			"good test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since the request can only be used within a gomega assertion,
			// register the fail handler, and a recover function
			gomega.RegisterFailHandler(ginkgo.Fail)
			defer ginkgo.GinkgoRecover()

			actual := NewResponse(tt.code, tt.body)
			expect := &Response{expectedCode: tt.code, expectedBody: tt.body, recorder: httptest.NewRecorder()}
			assert.Equal(t, expect, actual)

			actual.recorder.WriteHeader(tt.code)
			actual.recorder.WriteString(tt.body)

			var actualData Data
			actual.Unmarshall(&actualData)
			assert.Equal(t, d1, actualData)
		})
	}
}
