package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// router.GET("/albums", handlers.getAlbums)
	// router.GET("/albums/:id", )
	// router.POST("/albums",)

	router.Run("localhost:8080")
}
