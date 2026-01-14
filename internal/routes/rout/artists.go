package routes

import (
	"mess/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func ArtistRoutes(r *gin.Engine, handler handlers.HandlerInterface) {
	artist := r.Group("/artists")
	{
		artist.POST("/create", handler.CreateArtist)
		artist.GET("/list", handler.ArtistsByName)
		artist.GET("/get/:id", handler.ArtistByID)
		artist.GET("/update/:id", handler.ArtistsByName)

		a := artist.GET("/:id")
		{
			a.GET("/list", handler.TracksByArtist)
		}
	}
}
