package routes

import (
	"mess/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func TracksRoutes(r *gin.Engine, handler handlers.HandlerInterface) {
	track := r.Group("/track")
	{
		track.POST("/create", handler.AddTrack)
		track.GET("/list", handler.TracksByTitle)
		track.GET("/:id", handler.TrackByID)
		track.GET("/play/:id", handler.Stream)
		track.POST("/upload", handler.UploadFile)
		track.GET("/up", handler.Upload)
	}
}
