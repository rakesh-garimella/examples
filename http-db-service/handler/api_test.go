package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

func TestSwaggerAPIHandler(t *testing.T) {
	// fake api file
	ioutil.WriteFile("api.yaml", []byte("API Specs"), 0644)
	defer func() { assert.NoError(t, os.Remove("api.yaml")) }()

	router := mux.NewRouter()
	router.HandleFunc("/api.yaml", SwaggerAPIHandler)
	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/api.yaml", ts.URL))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
}

func TestSwaggerAPIRedirectHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/", SwaggerAPIRedirectHandler)
	router.HandleFunc("/api.yaml", func(w http.ResponseWriter, r *http.Request) {
		// just send accepted response to check that the redirect / -> /api.yaml happened.
		w.WriteHeader(http.StatusAccepted)
	})
	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
}
