package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreas-kokkalis/dock_server/pkg/api/store"
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

	r := httptest.NewRequest("POS", "http://localhost/foo", nil)
	w := httptest.NewRecorder()
	handler(w, r, httprouter.Params{})
	assert.Equal(t, 200, w.Code)

	/*response := w.Result()
	var p []byte
	response.Body.Read(p)
	log.Println(&p)
	*/
}

func TestAdminLogout(t *testing.T) {
	t.Parallel()

	redisRepo := store.NewRedisRepo(redismock.NewRedisMock().WithDel(0, nil))
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
	redisRepo = store.NewRedisRepo(redismock.NewRedisMock().WithDel(-1, errors.New("error")))
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
	redisRepo := store.NewRedisRepo(redismock.NewRedisMock().WithExists(true, nil))
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
	redisRepo = store.NewRedisRepo(redismock.NewRedisMock().WithExists(false, errors.New("error")))
	handler = SessionAuth(NewService(nil, redisRepo), h)
	http.SetCookie(w, &http.Cookie{Name: "ses", Value: "adm:7ff10abb653dead4186089acbd2b7891"})
	r = &http.Request{Header: http.Header{"Cookie": w.HeaderMap["Set-Cookie"]}}
	handler(w, r, httprouter.Params{})
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}
