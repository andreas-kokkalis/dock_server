package image

// var configDir = "../../../conf/"

/*
func TestListImages(t *testing.T) {
	assert := assert.New(t)
	// Initialize the configuration manager
	var c config.Config
	c, err = config.NewConfig(configDir, vars.Mode)
	if err != nil {
		log.Fatal(err)
	}

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
*/
