package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories"
	"github.com/andreas-kokkalis/dock_server/pkg/api/repositories/repomocks"
	"github.com/andreas-kokkalis/dock_server/pkg/drivers/redis/redismock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	t.Parallel()
	s := NewService(nil, nil)
	assert.NotNil(t, s)
}

func TestOauth(t *testing.T) {

	s := NewService(nil, nil)
	h := func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.WriteHeader(401)
	}
	handler := OAuth(s, h)

	r := httptest.NewRequest("POST", "http://localhost/foo", nil)
	w := httptest.NewRecorder()
	handler(w, r, httprouter.Params{})
	assert.Equal(t, 200, w.Code)
}

func TestAdminLogout(t *testing.T) {
	t.Parallel()

	redisRepo := repositories.NewRedisRepo(redismock.NewRedisMock().WithDel(0, nil))
	s := NewService(nil, redisRepo)
	handler := AdminLogout(s)
	r := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{Name: "ses", Value: "1"})
	r = &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	redisRepo = repositories.NewRedisRepo(redismock.NewRedisMock().WithDel(-1, errors.New("error")))
	handler = AdminLogout(NewService(nil, redisRepo))
	http.SetCookie(w, &http.Cookie{Name: "ses", Value: "1"})
	r = &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSessionAuth(t *testing.T) {
	t.Parallel()

	h := func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		w.WriteHeader(200)
		return
	}
	redisRepo := repositories.NewRedisRepo(redismock.NewRedisMock().WithExists(true, nil))
	s := NewService(nil, redisRepo)
	handler := SessionAuth(s, h)
	r := httptest.NewRequest("GET", "/foo", nil)
	w := httptest.NewRecorder()
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{Name: "ses", Value: "1"})
	r = &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// tests that the mock function h returns StatusOK
	w = httptest.NewRecorder()
	http.SetCookie(w, &http.Cookie{Name: "ses", Value: "adm:7ff10abb653dead4186089acbd2b7891"})
	r = &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	redisRepo = repositories.NewRedisRepo(redismock.NewRedisMock().WithExists(false, errors.New("error")))
	handler = SessionAuth(NewService(nil, redisRepo), h)
	http.SetCookie(w, &http.Cookie{Name: "ses", Value: "adm:7ff10abb653dead4186089acbd2b7891"})
	r = &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func TestAdminLogin(t *testing.T) {
	admin := api.Admin{
		ID:       1,
		Username: "foo",
		Password: "$2a$10$4F5Hpu0NM8Uy4bI/XQWKDO552uK77WwNpi3zIforzLngziZVszk06",
	}
	redisRepo := repositories.NewRedisRepo(redismock.NewRedisMock().WithExists(true, nil))
	adminMatches := repomocks.NewAdminDBRepositoryMock().WithGetAdminByUsername(admin, nil)

	tests := []struct {
		service    Service
		request    *http.Request
		expectCode int
		name       string
	}{
		{
			service: NewService(
				adminMatches,
				redisRepo,
			),
			request:    httptest.NewRequest(http.MethodGet, "/", nil),
			expectCode: http.StatusBadRequest,
			name:       "No JSON Body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := AdminLogin(tt.service)
			w := httptest.NewRecorder()
			handler(w, tt.request, nil)
			assert.Equal(t, tt.expectCode, w.Code, tt.name)
		})

	}

}
