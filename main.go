package main

import (
	"github.com/gin-gonic/gin"
	"github.com/renu-ramesh/robot-apocalypse-docker/handlers"
)

func main() {
	h := handlers.NewHandler()
	router := gin.Default()
	router.GET("/albums", h.getAlbums)
	router.GET("/albums/:id", h.getAlbumByID)
	router.POST("/albums", h.postAlbums)

	router.Run("localhost:8080")
}
