package services

import (
	"context"
	"mess/internal/model"
)

type artistService interface {
	CreateArtist(artist *model.Artist, ctx context.Context) error
	ArtistByID(id int, ctx context.Context) (*model.Artist, error)
	Artists(ctx context.Context) ([]model.Artist, error)
	ArtistsByName(name string, ctx context.Context) ([]model.Artist, error)
	FindArtists(input string) (error, []model.Artist)
}

// ArtistService creates a new artist
func (m *Service) CreateArtist(artist *model.Artist, ctx context.Context) error {
	return m.Repo.CreateArtist(artist, ctx)
}

// ArtistService retrieves an artist by ID
func (m *Service) ArtistByID(id int, ctx context.Context) (*model.Artist, error) {
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

// FindArtists retrieves artists by input string.
func (m *Service) FindArtists(input string) (error, []model.Artist) {
	return m.Repo.FindArtists(input)
}
