package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type JsonResponse struct {
	Hostname string `json:"hostname"`
	Count    int    `json:"count"`
}

func TestRoot(t *testing.T) {
	r := setupRouter("", "")
	testSrv := httptest.NewServer(r)
	defer testSrv.Close()

	resp, err := http.Get(testSrv.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	byteArray, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	b := string(byteArray)
	assert.Equal(t, "Hi, GET /incr to check session.", b)
}

func TestCookie(t *testing.T) {
	r := setupRouter("", "")
	testSrv := httptest.NewServer(r)
	defer testSrv.Close()

	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)
	client := &http.Client{Jar: jar}

	for i := 0; i < 3; i++ {
		resp, err := client.Get(testSrv.URL + "/incr")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var j JsonResponse
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&j)
		assert.NoError(t, err)
		assert.Equal(t, i, j.Count)
	}
}

func TestRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	r := setupRouter("redis:6379", "")
	testSrv := httptest.NewServer(r)
	defer testSrv.Close()

	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)
	client := &http.Client{Jar: jar}

	for i := 0; i < 3; i++ {
		resp, err := client.Get(testSrv.URL + "/incr")
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var j JsonResponse
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&j)
		assert.NoError(t, err)
		assert.Equal(t, i, j.Count)
	}
}
