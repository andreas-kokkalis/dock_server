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

			// w := httptest.NewRecorder()
			expect := &Response{expectedCode: tt.code, expectedBody: tt.body, recorder: httptest.NewRecorder()}
			assert.Equal(t, expect, actual)

			actual.recorder.WriteHeader(tt.code)
			actual.recorder.WriteString(tt.body)

			// Test the toString log func
			// js, err := json.MarshalIndent(actual, "", " ")
			// assert.NoError(t, err)
			// actualRequestString := actual.toString()
			// expectedRequestString := string(js)
			// assert.Contains(t, actualRequestString, expectedRequestString)
		})
	}
}
