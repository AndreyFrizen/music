package services

import "mess/internal/model"

type albumService interface {
	AddAlbum(album *model.Album) error
	AlbumByID(id string) (*model.Album, error)
}

// AddAlbum adds a new album
func (m *Service) AddAlbum(album *model.Album) error {
	return m.Repo.AddAlbum(album)
}

// AlbumByID retrieves an album by its ID.
func (m *Service) AlbumByID(id string) (*model.Album, error) {
	return m.Repo.AlbumByID(id)
}
