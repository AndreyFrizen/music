package auth

import (
	"mess/internal/model"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func UserService(user model.User, store model.UserRepository) error {
	err := validate.Struct(user)
	err = user.EncryptPassword()
	err = store.CreateUser(&user)

	if err != nil {
		return err
	}

	return nil
}

func ArtistService(artist model.Artist, store model.ArtistRepository) error {
	err := store.AddArtist(&artist)

	if err != nil {
		return err
	}

	return nil
}

func TrackService(track model.Track, store model.TrackRepository) error {
	err := store.AddTrack(&track)

	if err != nil {
		return err
	}

	return nil
}

func AlbumService(album model.Album, store model.AlbumRepository) error {
	err := store.AddAlbum(&album)

	if err != nil {
		return err
	}

	return nil
}

func PlaylistService(playlist model.Playlist, store model.PlaylistRepository) error {
	err := store.CreatePlaylist(&playlist)

	if err != nil {
		return err
	}

	return nil
}

func PlaylistTrackService(track model.PlaylistTrack, store model.PlaylistRepository) error {
	err := store.AddTrackToPlaylist(&track)

	if err != nil {
		return err
	}

	return nil
}
