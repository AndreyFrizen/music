package handlers

import (
	"net/http"
	"playlist-service/internal/domain/model"
	"playlist-service/proto/playlist"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler struct {
	playlist.UnimplementedPlaylistServiceServer
	service PlaylistService
}

type PlaylistService interface {
	CreatePlaylist(*model.Playlist, *gin.Context) error
	DeletePlaylist(int64, *gin.Context) error
	PlaylistByID(int64, *gin.Context) (*model.Playlist, error)
	UpdatePlaylist(int64, *model.Playlist, *gin.Context) error
	AddTrack(int64, int64, *gin.Context) error
	RemoveTrack(int64, int64, *gin.Context) error
}
