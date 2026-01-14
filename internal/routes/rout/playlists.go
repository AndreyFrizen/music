package routes

import (
	"mess/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func PlaylistsRoutes(r *gin.Engine, handler handlers.HandlerInterface) {
	playlists := r.Group("/playlist")
	{
		playlists.POST("/create", handler.CreatePlaylist)
		playlists.GET("/:id", handler.PlaylistByID)
		playlists.DELETE("/update/:id", handler.DeletePlaylist)

		p := playlists.GET("/:id")
		{
			p.POST("/add", handler.AddTrackToPlaylist)
			p.POST("/delete", handler.DeleteTrackFromPlaylist)
			p.GET("/:id", handler.TrackFromPlaylist)
		}
	}
}
