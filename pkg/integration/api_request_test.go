package integration

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		method, url string
		body        interface{}
		session     string
		name        string
	}{
		{
			http.MethodGet,
			"/foo/bar",
			nil,
			"",
			"empty body",
		},
		{
			http.MethodGet,
			"/foo/bar",
			api.Img{ID: "abc", RepoTags: []string{"foo"}, CreatedAt: time.Date(2017, 12, 20, 10, 24, 0, 0, time.UTC)},
			"",
			"with body",
		},
		{
			http.MethodGet,
			"/foo/bar",
			api.Img{ID: "abc", RepoTags: []string{"foo"}, CreatedAt: time.Date(2017, 12, 20, 10, 24, 0, 0, time.UTC)},
			"val",
			"with body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Since the request can only be used within a gomega assertion,
			// register the fail handler, and a recover function
			gomega.RegisterFailHandler(ginkgo.Fail)
			defer ginkgo.GinkgoRecover()

			actual := NewRequest(tt.method, tt.url, tt.body)

			if tt.session != "" {
				actual.WithSessionCookie(tt.session)
			}

			jsonBody, err := json.Marshal(tt.body)
			assert.NoError(t, err)
			r := httptest.NewRequest(tt.method, tt.url, ioutil.NopCloser(bytes.NewReader(jsonBody)))
			if tt.session != "" {
				r.AddCookie(&http.Cookie{Name: "ses", Value: tt.session})
			}
			expect := &Request{method: tt.method, url: tt.url, body: tt.body, HTTPRequest: r}
			assert.Equal(t, expect, actual)

			// Test the toString log func
			js, err := json.MarshalIndent(LogRequest{Req{tt.method, tt.url, tt.body}}, "", "  ")
			assert.NoError(t, err)
			actualRequestString := actual.pretty()
			expectedRequestString := string(js)
			assert.Equal(t, actualRequestString, expectedRequestString)
		})
	}
}
