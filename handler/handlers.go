package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wie.gg/lnk/shortener"
	"wie.gg/lnk/store"
)

// Request model definition
type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var req UrlCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
		return
	}

	shortUrl := shortener.GenerateShortLink(req.LongUrl, req.UserId)
	store.SaveUrlMapping(shortUrl, req.LongUrl, req.UserId)

	host := "http://localhost:8080/"
	c.JSON(200, gin.H{
		"message": "URL successfully created",
		"url": host + shortUrl,
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initalUrl := store.RetrieveOriginalUrl(shortUrl)

	c.Redirect(302, initalUrl)
}