package handlers

import (
	services "mess/internal/service"
)

type Handler struct {
	service services.InterfaceService
}

type HandlerInterface interface {
	userHandler
	artistHandler
	albumHandler
	playlistHandler
	trackHandler
}

func NewHandler(repo services.InterfaceService) *Handler {
	return &Handler{
		service: repo,
	}
}
