package handler

import (
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

var _env = "../.env"

func init() {
	_store, _mock = redismock.NewClientMock()
	store.InitializeStore(_store)
}

func TestGet(t *testing.T) {
	router := SetupRouter(&_env)

	_mock.ExpectGet("test").SetVal("http://google.com")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
}

func TestAuthNotProvided(t *testing.T) {
	router := SetupRouter(&_env)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}
