package repository

import (
	"collection-service/internal/app/database"
	"collection-service/internal/domain/errors"
	"collection-service/internal/domain/models"
	"context"
	"strconv"
)

type store struct {
	db *database.DB
}

func NewRepository(db *database.DB) *store {
	return &store{db: db}
}

func (r *store) AddAlbum(ctx context.Context, album *models.Album) (int64, error) {
	const op = "repository.CollectionRepository.AddAlbum"

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO albums (user_id, album_id)
		VALUES ($1, $2)
	`, album.UserId, album.AlbumId)
	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	go r.setAlbumToCache(ctx, strconv.Itoa(int(album.UserId)), album)

	return album.AlbumId, nil
}

func (r *store) GetAlbums(ctx context.Context, userId int64) ([]*models.Album, error) {
	const op = "repository.CollectionRepository.GetAlbum"

	var result []*models.Album
	err := r.db.QueryRowContext(ctx, `
		SELECT *
		FROM albums
		WHERE user_id = $1
	`, userId).Scan(&result)
	if err != nil {
		return nil, errors.DatabaseError(op, err)
	}

	return result, nil
}

func (r *store) DeleteAlbum(ctx context.Context, userId int64, albumId int64) error {
	const op = "repository.CollectionRepository.DeleteAlbum"

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM albums
		WHERE user_id = $1 AND album_id = $2
	`, userId, albumId)
	if err != nil {
		return errors.DatabaseError(op, err)
	}

	return nil
}

func (r *store) AddArtist(ctx context.Context, artist *models.Artist) (int64, error) {
	const op = "repository.CollectionRepository.AddArtist"

	var artistId int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO artists (user_id, artist_id)
		VALUES ($1, $2)
		RETURNING artist_id
	`, artist.UserId, artist.ArtistId).Scan(&artistId)
	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	return artistId, nil
}

func (r *store) GetArtists(ctx context.Context, userId int64) ([]*models.Artist, error) {
	const op = "repository.CollectionRepository.GetArtists"

	var result []*models.Artist
	err := r.db.QueryRowContext(ctx, `
		SELECT *
		FROM artists
		WHERE user_id = $1
	`, userId).Scan(&result)
	if err != nil {
		return nil, errors.DatabaseError(op, err)
	}

	return result, nil
}

func (r *store) DeleteArtist(ctx context.Context, userId int64, artistId int64) error {
	const op = "repository.CollectionRepository.DeleteArtist"

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM artists
		WHERE user_id = $1 AND artist_id = $2
	`, userId, artistId)
	if err != nil {
		return errors.DatabaseError(op, err)
	}

	return nil
}

func (r *store) AddTrack(ctx context.Context, track *models.Track) (int64, error) {
	const op = "repository.CollectionRepository.AddTrack"

	var trackId int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO tracks (user_id, track_id)
		VALUES ($1, $2)
		RETURNING track_id
	`, track.UserId, track.TrackId).Scan(&trackId)
	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	return trackId, nil
}

func (r *store) GetTracks(ctx context.Context, userId int64) ([]*models.Track, error) {
	const op = "repository.CollectionRepository.GetTracks"

	var result []*models.Track
	err := r.db.QueryRowContext(ctx, `
		SELECT *
		FROM tracks
		WHERE user_id = $1
	`, userId).Scan(&result)
	if err != nil {
		return nil, errors.DatabaseError(op, err)
	}

	return result, nil
}

func (r *store) DeleteTrack(ctx context.Context, userId int64, trackId int64) error {
	const op = "repository.CollectionRepository.DeleteTrack"

	_, err := r.db.ExecContext(ctx, `
		DELETE FROM tracks
		WHERE user_id = $1 AND track_id = $2
	`, userId, trackId)
	if err != nil {
		return errors.DatabaseError(op, err)
	}

	return nil
}
