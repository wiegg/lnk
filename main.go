package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"wie.gg/lnk/handler"
	"wie.gg/lnk/store"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "HEYYYY",
		})
	})

	r.POST("/", func (c *gin.Context)  {
		handler.CreateShortUrl(c)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	store.InitializeStore(nil)

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}