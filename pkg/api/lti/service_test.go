package lti

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/portmapper"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories/repomocks"
	"github.com/stretchr/testify/assert"

	"github.com/julienschmidt/httprouter"
)

type request struct {
	form   url.Values
	header http.Header
	method string
	target string
	body   io.Reader
}

func newRequest(method, target string, body io.Reader) *request {
	return &request{
		method: method,
		target: target,
		body:   body,
	}
}

func (r *request) withPostForm(kv map[string]string) *request {
	form := url.Values{}
	for k, v := range kv {
		form.Set(k, v)
	}
	r.form = form
	return r
}

func (r *request) withHeader() *request {
	w := httptest.NewRecorder()
	r.header = http.Header{"Content-Type": w.HeaderMap["application/x-www-form-urlencoded"]}
	return r
}
func (r *request) do() *http.Request {
	// req := httptest.NewRequest(r.method, r.target, r.body)
	req := &http.Request{
		Method:     r.method,
		RequestURI: r.target,
		PostForm:   r.form,
		Header:     r.header,
	}
	if r.body != nil {
		switch v := r.body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
		default:
			req.ContentLength = -1
		}
		if rc, ok := r.body.(io.ReadCloser); ok {
			req.Body = rc
		} else {
			req.Body = ioutil.NopCloser(r.body)
		}
	}
	return req
}

var tmplPath = "../../../"

func TestLaunch(t *testing.T) {

	tests := []struct {
		service        Service
		request        *http.Request
		params         []httprouter.Param
		expectCode     int
		expectResponse bool
		expectResp     Resp
		name           string
	}{
		{
			service:        NewService(nil, nil, nil, "../../"),
			request:        httptest.NewRequest(http.MethodPost, "/", nil),
			params:         nil,
			expectCode:     http.StatusInternalServerError,
			expectResponse: false,
			name:           "invalid template path",
		},
		{
			service:        NewService(nil, nil, nil, tmplPath),
			request:        httptest.NewRequest(http.MethodPost, "/", nil),
			params:         []httprouter.Param{{Key: "id", Value: "foo"}},
			expectCode:     http.StatusBadRequest,
			expectResponse: true,
			expectResp:     Resp{Error: ErrInvalidImageID},
			name:           "invalid imageID",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock(),
				repomocks.NewDockerRepositoryMock(),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 5),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", nil).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusBadRequest,
			expectResponse: true,
			expectResp:     Resp{Error: ErrInvalidFormData},
			name:           "invalid form data",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock().
					WithUserRunConfigExists(false, errors.New("redis error")),
				repomocks.NewDockerRepositoryMock(),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 5),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(""))).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusInternalServerError,
			expectResponse: true,
			expectResp:     Resp{Error: api.ErrServerError},
			name:           "with redis session exists error",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock().
					WithUserRunConfigExists(true, nil).
					WithUserRunConfigGet(api.RunConfig{}, errors.New("redis error")),
				repomocks.NewDockerRepositoryMock(),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 5),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(""))).withPostForm(map[string]string{"user_id": "1"}).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusInternalServerError,
			expectResponse: true,
			expectResp:     Resp{Error: api.ErrServerError},
			name:           "with redis session exists - get key errors",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock().
					WithUserRunConfigExists(true, nil).
					WithUserRunConfigGet(api.RunConfig{}, nil).
					WithUserRunConfigSet(errors.New("redis error")),
				repomocks.NewDockerRepositoryMock(),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 5),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(""))).withPostForm(map[string]string{"user_id": "1"}).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusInternalServerError,
			expectResponse: true,
			expectResp:     Resp{Error: api.ErrServerError},
			name:           "with redis session exists - set key errors",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock().
					WithUserRunConfigExists(false, nil),
				repomocks.NewDockerRepositoryMock(),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 0),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(""))).withPostForm(map[string]string{"user_id": "1"}).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusInternalServerError,
			expectResponse: true,
			expectResp:     Resp{Error: ErrResourceQuotaExceeded},
			name:           "with redis session not exists - mapper errors",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock().
					WithUserRunConfigExists(false, nil),
				repomocks.NewDockerRepositoryMock().
					WithContainerRun(api.RunConfig{}, errors.New("container run error")),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 2),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(""))).withPostForm(map[string]string{"user_id": "1"}).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusInternalServerError,
			expectResponse: true,
			expectResp:     Resp{Error: ErrContainerRun},
			name:           "with redis session not exists - container run errors",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock().
					WithUserRunConfigExists(false, nil).
					WithUserRunConfigSet(errors.New("redis error set")),
				repomocks.NewDockerRepositoryMock().
					WithContainerRun(api.RunConfig{}, nil),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 2),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(""))).withPostForm(map[string]string{"user_id": "1"}).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusInternalServerError,
			expectResponse: true,
			expectResp:     Resp{Error: api.ErrServerError},
			name:           "with redis session not exists - setting run config errors",
		},
		{
			service: NewService(
				repomocks.NewRedisRepositoryMock().
					WithUserRunConfigExists(false, nil).
					WithUserRunConfigSet(nil).
					WithUserRunKeyGet("key"),
				repomocks.NewDockerRepositoryMock().
					WithContainerRun(api.RunConfig{}, nil),
				portmapper.NewPortMapper(repomocks.NewRedisRepositoryMock(), 2),
				tmplPath),
			request:        newRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(""))).withPostForm(map[string]string{"user_id": "1"}).do(),
			params:         []httprouter.Param{{Key: "id", Value: "83364c85cafc"}},
			expectCode:     http.StatusOK,
			expectResponse: true,
			expectResp:     getResp(api.RunConfig{}),
			name:           "with redis session not exists - No errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.service.Launch(w, tt.request, tt.params)

			assert.Equal(t, tt.expectCode, w.Code, tt.name)

			if tt.expectResponse {
				tmp := path.Join(tmplPath, "templates/html/assignment.html")
				templ, err := template.ParseFiles(tmp)
				assert.NoError(t, err)
				wr := httptest.NewRecorder()
				wr.WriteHeader(tt.expectCode)
				assert.NoError(t, templ.Execute(wr, tt.expectResp))
				assert.Equal(t, wr.Body.String(), w.Body.String())
			}

			if w.Code == http.StatusOK {
				assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
			}
		})
	}
}
