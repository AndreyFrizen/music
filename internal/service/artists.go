package services

import (
	"context"
	"mess/internal/model"
)

type artistService interface {
	CreateArtist(artist *model.Artist, ctx context.Context) error
	ArtistByID(id string, ctx context.Context) (*model.Artist, error)
	Artists(ctx context.Context) ([]model.Artist, error)
	ArtistsByName(name string, ctx context.Context) ([]model.Artist, error)
}

// ArtistService creates a new artist
func (m *Service) CreateArtist(artist *model.Artist, ctx context.Context) error {
	return m.Repo.CreateArtist(artist, ctx)
}

// ArtistService retrieves an artist by ID
func (m *Service) ArtistByID(id string, ctx context.Context) (*model.Artist, error) {
	return m.Repo.ArtistByID(id, ctx)
}

// ArtistService retrieves all artists
func (m *Service) Artists(ctx context.Context) ([]model.Artist, error) {
	return m.Repo.Artists(ctx)
}

// ArtistService retrieves artists by name
func (m *Service) ArtistsByName(name string, ctx context.Context) ([]model.Artist, error) {
	return m.Repo.ArtistsByName(name, ctx)
}
