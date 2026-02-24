package routes

import (
	"mess/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func ArtistRoutes(r *gin.Engine, handler handlers.HandlerInterface) {
	artist := r.Group("/artists")
	{
		artist.GET("/ws", handler.ArtistWebSocket)
		artist.GET("/find", handler.FindArtists)
		artist.POST("/create", handler.CreateArtist)
		artist.GET("/list", handler.Artists)
		artist.GET("/:id", handler.ArtistByID)
		artist.GET("/name/:name", handler.ArtistsByName)
		artist.GET("/tracks/:id", handler.TracksByArtist)
	}
}
