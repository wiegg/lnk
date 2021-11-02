package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"wie.gg/lnk/store"
)

var _store *redis.Client
var _mock redismock.ClientMock

func init() {
	_store, _mock = redismock.NewClientMock()
	store.InitializeStore(_store)
}

func TestGet(t *testing.T) {
	router := SetupRouter()

	_mock.ExpectGet("test").SetVal("http://google.com")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
}

func TestGetJson(t *testing.T) {
	router := SetupRouter()

	_mock.ExpectGet("test").SetVal("http://google.com")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?json=true", nil)
	router.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Result().Body)

	type AssertResponse struct {
		InitialUrl string `json:"initialUrl"`
	}

	var test AssertResponse
	json.Unmarshal(body, &test)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "http://google.com", test.InitialUrl)

}

func TestHealth(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuthNotProvided(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}
