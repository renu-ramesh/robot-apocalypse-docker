package main

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/renu-ramesh/robot-apocalypse-docker/handlers"
)

func main() {
	var h handlers.Handler
	router := gin.Default()
	router.GET("/albums", h.getAlbums)
	// router.GET("/albums/:id", h.getAlbumByID)
	// router.POST("/albums", h.postAlbums)

	router.Run("localhost:8080")
}
