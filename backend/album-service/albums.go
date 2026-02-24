package routes

import (
	"mess/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func AlbumsRoutes(r *gin.Engine, handler handlers.HandlerInterface) {
	albums := r.Group("/albums")
	{
		albums.POST("/create", handler.AddAlbum)
		albums.GET("/:id", handler.AlbumByID)
		albums.GET("/title/:title", handler.AlbumsByTitle)
	}
}
