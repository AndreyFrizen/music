package services

import (
	"context"
	"mess/internal/model"
)

type albumService interface {
	AddAlbum(album *model.Album, ctx context.Context) error
	AlbumByID(id int, ctx context.Context) (*model.Album, error)
	AlbumsByArtist(artistID int, ctx context.Context) ([]model.Album, error)
	AlbumsByTitle(title string, ctx context.Context) ([]model.Album, error)
}

// AddAlbum adds a new album
func (m *Service) AddAlbum(album *model.Album, ctx context.Context) error {
	if err := m.Repo.AddAlbum(album, ctx); err != nil {
		return err
	}
	return nil
}

// AlbumByID retrieves an album by its ID.
func (m *Service) AlbumByID(id int, ctx context.Context) (*model.Album, error) {
	return m.Repo.AlbumByID(id, ctx)
}

// AlbumsByArtist retrieves albums by artist ID.
func (m *Service) AlbumsByArtist(artistID int, ctx context.Context) ([]model.Album, error) {
	return m.Repo.AlbumsByArtistID(artistID, ctx)
}

// AlbumsByArtistName retrieves albums by artist name.
func (m *Service) AlbumsByTitle(title string, ctx context.Context) ([]model.Album, error) {
	return m.Repo.AlbumsByTitle(title, ctx)
}
