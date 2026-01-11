package services

import (
	"mess/internal/repository/store"

	"github.com/go-playground/validator/v10"
)

type InterfaceService interface {
	userService
	artistService
	trackService
	albumService
	playlistService
}

type Service struct {
	Repo store.Repository
}

func NewService(repo store.Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}
