package handler

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"wie.gg/lnk/middleware"
	"wie.gg/lnk/shortener"
	"wie.gg/lnk/store"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Authorization"},
	}))

	r.POST("/", middleware.EnsureValidToken(), CreateShortUrl)

	r.GET("/health", HandleHealth)

	r.GET("/:shortUrl", HandleShortUrlRedirect)

	return r
}

// Request model definition
type UrlCreationRequest struct {
	LongUrl string `json:"url" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	t := c.Request.Context().Value(jwtmiddleware.ContextKey{})

	if t == nil {
		return
	}

	claims := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*middleware.CustomClaims)
	if !claims.HasScope("create:link") {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Insufficient permissions",
		})
	}

	var req UrlCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := shortener.GenerateShortLink(req.LongUrl, claims.User.Id)
	store.SaveUrlMapping(shortUrl, req.LongUrl, claims.User.Id)

	c.JSON(200, gin.H{
		"message": "URL successfully created",
		"url":     "/" + shortUrl,
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	initalUrl := store.RetrieveOriginalUrl(shortUrl)

	asJson := c.DefaultQuery("json", "false")

	if asJson == "false" {
		c.Redirect(302, initalUrl)
	} else {
		c.JSON(200, gin.H{
			"shortUrl":   shortUrl,
			"initialUrl": initalUrl,
		})
	}
}

func HandleHealth(c *gin.Context) {
	c.Status(200)
}
