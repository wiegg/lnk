package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testService = &StorageService{}

func init() {
	testService = InitializeStore()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testService.redis != nil)
}

func TestInsertAndRetrieval(t *testing.T) {
	initialLink := "https://www.guru3d.com/news-story/spotted-ryzen-threadripper-pro-3995wx-processor-with-8-channel-ddr4,2.html"
	userUUId := "e0dba740-fc4b-4977-872c-d360239e6b1a"
	shortURL := "Jsz4k57oAX"

	// Persist data mapping
	SaveUrlMapping(shortURL, initialLink, userUUId)

	// Retrieve initial URL
	retrievedUrl := RetrieveOriginalUrl(shortURL)

	assert.Equal(t, initialLink, retrievedUrl)
}