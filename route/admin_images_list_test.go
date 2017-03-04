package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreas-kokkalis/dock-server/srv"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestListImages(t *testing.T) {
	assert := assert.New(t)
	srv.InitDep("../conf/")

	router := httprouter.New()
	router.GET("/v0/admin/images", ListImages)
	r, err := http.NewRequest("GET", "/v0/admin/images", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	assert.Equal(w.Code, http.StatusOK, "It should be equal")

	decoder := json.NewDecoder(w.Body)
	var data Response
	err = decoder.Decode(&data)
	assert.Equal(nil, err, "It should not return an error")

	resp := w.Result()
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	t.Logf("%s", bodyString)

	fmt.Printf("LIST_IMAGES:\t %d - %s", w.Code, w.Body.String())

}
