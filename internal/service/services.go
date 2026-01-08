package services

import (
	"mess/internal/model"
	"mess/internal/repository/store"

	"github.com/go-playground/validator/v10"
)

type MusicService struct {
	repo store.Repository
}

type MusicServices interface {
	RegisterService(*model.User) error
	ArtistService(*model.Artist) error
	TrackService(*model.Track) error
	AlbumService(*model.Album) error
	PlaylistService(*model.Playlist) error
	PlaylistTrackService(*model.PlaylistTrack) error
}

func NewMusicService(repo store.Repository) *MusicService {
	return &MusicService{
		repo: repo,
	}
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (m *MusicService) RegisterService(user *model.User) error {
	err := validate.Struct(user)
	err = user.EncryptPassword()
	err = m.repo.CreateUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (m *MusicService) ArtistService(artist *model.Artist) error {
	err := m.repo.AddArtist(artist)

	if err != nil {
		return err
	}

	return nil
}

func (m *MusicService) TrackService(track *model.Track) error {
	err := m.repo.AddTrack(track)

	if err != nil {
		return err
	}

	return nil
}

func (m *MusicService) AlbumService(album *model.Album) error {
	err := m.repo.AddAlbum(album)

	if err != nil {
		return err
	}

	return nil
}

func (m *MusicService) PlaylistService(playlist *model.Playlist) error {
	err := m.repo.CreatePlaylist(playlist)

	if err != nil {
		return err
	}

	return nil
}

func (m *MusicService) PlaylistTrackService(track *model.PlaylistTrack) error {
	err := m.repo.AddTrackToPlaylist(track)

	if err != nil {
		return err
	}

	return nil
}

// func (m *MusicService) GetTrackStream() error {
// 	_, err := os.Open("~/project/music/static/1.mp3")
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
