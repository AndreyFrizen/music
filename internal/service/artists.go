package services

import "mess/internal/model"

type artistService interface {
	CreateArtist(artist *model.Artist) error
	ArtistByID(id string) (*model.Artist, error)
	Artists() ([]model.Artist, error)
	ArtistsByName(name string) ([]model.Artist, error)
}

// ArtistService creates a new artist
func (m *Service) CreateArtist(artist *model.Artist) error {
	return m.Repo.CreateArtist(artist)
}

// ArtistService retrieves an artist by ID
func (m *Service) ArtistByID(id string) (*model.Artist, error) {
	return m.Repo.ArtistByID(id)
}

// ArtistService retrieves all artists
func (m *Service) Artists() ([]model.Artist, error) {
	return m.Repo.Artists()
}

// ArtistService retrieves artists by name
func (m *Service) ArtistsByName(name string) ([]model.Artist, error) {
	return m.Repo.ArtistsByName(name)
}
