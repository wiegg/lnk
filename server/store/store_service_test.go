package store

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

var testService = &StorageService{}
var mock redismock.ClientMock
var db *redis.Client

func init() {
	db, mock = redismock.NewClientMock()
	testService = InitializeStore(db)
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testService.redis != nil)
}

func TestInsertAndRetrieval(t *testing.T) {
	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"

	mock.ExpectSet(shortURL, initialLink, 6 * time.Hour).SetVal("")
	mock.ExpectGet(shortURL).SetVal(initialLink)

	// Persist data mapping
	SaveUrlMapping(shortURL, initialLink, userUUId)

	// Retrieve initial URL
	retrievedUrl := RetrieveOriginalUrl(shortURL)

	assert.Equal(t, initialLink, retrievedUrl)
}