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
		playlists.DELETE("/:id", handler.DeletePlaylist)
		playlists.POST("/add/:id", handler.AddTrackToPlaylist)
		playlists.DELETE("/delete/:id", handler.DeleteTrackFromPlaylist)
		playlists.GET("/track/:id", handler.TrackFromPlaylist)
	}
}
