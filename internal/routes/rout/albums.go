package routes

import (
	"mess/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func AlbumsRoutes(r *gin.Engine, handler handlers.HandlerInterface) {
	albums := r.Group("/albums")
	{
		albums.POST("/create", handler.AddAlbum)
		albums.GET("/list", handler.AlbumByID)
		albums.GET("/list/:id", handler.AlbumsByArtist)
		albums.GET("/list/:id", handler.AlbumsByTitle)
	}
}
