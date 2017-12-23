package auth

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories/repomocks"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/postgres"
	"github.com/jordic/lti"
	"github.com/jordic/lti/oauth"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	s := NewService(nil, nil)
	assert.NotNil(t, s)
}

func oauthRequest(uriPath string) *http.Request {
	u, _ := url.Parse(uriPath)
	p := &lti.Provider{
		Secret:      oauthSecret,
		URL:         uriPath,
		ConsumerKey: oauthKey,
		Method:      http.MethodPost,
	}
	p.SetSigner(oauth.GetHMACSigner(oauthSecret, ""))
	_, _ = p.Sign()
	return &http.Request{
		Method: http.MethodPost,
		Host:   u.Host,
		URL:    u,
		Body:   nil,
		Form:   p.Params(),
	}
}

func TestOauth(t *testing.T) {
	h := func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) { w.WriteHeader(200) }
	tests := []struct {
		service    Service
		request    *http.Request
		expectCode int
		name       string
	}{
		{
			service:    NewService(nil, nil),
			request:    httptest.NewRequest(http.MethodPost, "/", nil),
			expectCode: http.StatusBadRequest,
			name:       "oauth invalid",
		},
		{
			service:    NewService(nil, nil),
			request:    oauthRequest("https://example.com/lti"),
			expectCode: http.StatusOK,
			name:       "oauth valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.service.OAuth(h)(w, tt.request, nil)
			assert.Equal(t, tt.expectCode, w.Code, tt.name)
		})
	}
}

func cookieRequest(val string) *http.Request {
	w := httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{Name: "ses", Value: val})
	return &http.Request{
		Header:     http.Header{"Cookie": w.HeaderMap["Set-Cookie"]},
		Method:     http.MethodGet,
		RequestURI: "/",
	}
}

func TestAdminLogout(t *testing.T) {
	tests := []struct {
		service      Service
		request      *http.Request
		expectCode   int
		expectCookie *http.Cookie
		name         string
	}{
		{
			service:      NewService(nil, nil),
			request:      httptest.NewRequest(http.MethodGet, "/", nil),
			expectCode:   http.StatusUnauthorized,
			expectCookie: nil,
			name:         "session does not exist",
		},
		{
			service: NewService(
				nil,
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionDelete(errors.New("redis network error"))),
			request:      cookieRequest("1"),
			expectCode:   http.StatusInternalServerError,
			expectCookie: nil,
			name:         "deleting session errors",
		},
		{
			service: NewService(
				nil,
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionDelete(nil)),
			request:      cookieRequest("1"),
			expectCode:   http.StatusOK,
			expectCookie: &http.Cookie{Name: "ses", Value: "", Path: "/", Expires: time.Now()},
			name:         "deleting session errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.service.AdminLogout(w, tt.request, nil)
			assert.Equal(t, tt.expectCode, w.Code, tt.name)
			if tt.expectCookie != nil {
				assert.Equal(t, tt.expectCookie.String(), w.Header().Get("Set-Cookie"), tt.name)
			}
		})
	}
}

func TestSessionAuth(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) { w.WriteHeader(200) }
	tests := []struct {
		service      Service
		request      *http.Request
		expectCode   int
		expectCookie *http.Cookie
		name         string
	}{
		{
			service:    NewService(nil, nil),
			request:    httptest.NewRequest(http.MethodGet, "/", nil),
			expectCode: http.StatusUnauthorized,
			name:       "session does not exist",
		},
		{
			service:    NewService(nil, nil),
			request:    cookieRequest("1"),
			expectCode: http.StatusUnauthorized,
			name:       "invalid cookie hash",
		},
		{
			service: NewService(
				nil,
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionExists(false, errors.New("redis server error"))),
			request:    cookieRequest("adm:7ff10abb653dead4186089acbd2b7891"),
			expectCode: http.StatusInternalServerError,
			name:       "redis error",
		},
		{
			service: NewService(
				nil,
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionExists(false, nil)),
			request:    cookieRequest("adm:7ff10abb653dead4186089acbd2b7891"),
			expectCode: http.StatusUnauthorized,
			name:       "stale cookie does not exist in session storage",
		},
		{
			service: NewService(
				nil,
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionExists(true, nil)),
			request:    cookieRequest("adm:7ff10abb653dead4186089acbd2b7891"),
			expectCode: http.StatusOK,
			name:       "session exists SUCCESS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.service.SessionAuth(h)(w, tt.request, nil)
			assert.Equal(t, tt.expectCode, w.Code, tt.name)
		})
	}
}

func body(json string) io.Reader {
	return bytes.NewBuffer([]byte(json))
}

func TestAdminLogin(t *testing.T) {
	admin := api.Admin{
		ID:       1,
		Username: "foo",
		Password: "$2a$10$4F5Hpu0NM8Uy4bI/XQWKDO552uK77WwNpi3zIforzLngziZVszk06",
	}
	bodyGood := `{"usename":"admin", "password": "kthtest"}`
	bodyMismatch := `{"usename":"admin", "password": "foo"}`
	tests := []struct {
		service    Service
		body       io.Reader
		request    *http.Request
		expectCode int
		name       string
	}{
		{
			service:    NewService(nil, nil),
			body:       nil,
			expectCode: http.StatusBadRequest,
			name:       "No JSON Body",
		},
		{
			service: NewService(
				repomocks.NewAdminDBRepositoryMock().
					WithGetAdminByUsername(api.Admin{}, postgres.ErrNoResult),
				nil),
			body:       body(bodyGood),
			expectCode: http.StatusUnauthorized,
			name:       "Username does not exist",
		},
		{
			service: NewService(
				repomocks.NewAdminDBRepositoryMock().
					WithGetAdminByUsername(api.Admin{}, errors.New("database error")),
				nil),
			body:       body(bodyGood),
			expectCode: http.StatusInternalServerError,
			name:       "Database error",
		},
		{
			service: NewService(
				repomocks.NewAdminDBRepositoryMock().
					WithGetAdminByUsername(admin, nil),
				nil),
			body:       body(bodyMismatch),
			expectCode: http.StatusUnauthorized,
			name:       "Password mismatch",
		},
		{
			service: NewService(
				repomocks.NewAdminDBRepositoryMock().
					WithGetAdminByUsername(admin, nil),
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionKeyCreate("adminkey").
					WithAdminSessionExists(false, errors.New("session exists errors"))),
			body:       body(bodyGood),
			expectCode: http.StatusInternalServerError,
			name:       "Session exists errors",
		},
		{
			service: NewService(
				repomocks.NewAdminDBRepositoryMock().
					WithGetAdminByUsername(admin, nil),
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionKeyCreate("adminkey").
					WithAdminSessionExists(false, nil).
					WithAdminSessionSet(errors.New("set session errors"))),
			body:       body(bodyGood),
			expectCode: http.StatusInternalServerError,
			name:       "Session does not exist, session set error",
		},
		{
			service: NewService(
				repomocks.NewAdminDBRepositoryMock().
					WithGetAdminByUsername(admin, nil),
				repomocks.NewRedisRepositoryMock().
					WithAdminSessionKeyCreate("adminkey").
					WithAdminSessionExists(true, nil).
					WithAdminSessionSet(errors.New("set session errors"))),
			body:       body(bodyGood),
			expectCode: http.StatusOK,
			name:       "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", tt.body)
			w := httptest.NewRecorder()
			tt.service.AdminLogin(w, r, nil)
			assert.Equal(t, tt.expectCode, w.Code, tt.name, tt.name)
		})
	}
}
