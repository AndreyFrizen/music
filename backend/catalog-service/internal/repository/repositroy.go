package repository

import (
	"catalog-service/internal/app/database"
	"catalog-service/internal/domain/errors"
	"catalog-service/internal/domain/model"
	"context"
	"database/sql"
	"strconv"
)

type store struct {
	db *database.DB
}

func NewRepository(db *database.DB) *store {
	return &store{
		db: db,
	}
}

// Add track to database
func (s *store) CreateAlbum(ctx context.Context, a *model.Album) (int64, error) {
	const op = "repository.CatalogRepository.CreateAlbum"

	query := "INSERT INTO albums(title, artist_id, release_date) VALUES ($1, $2, $3) RETURNING id"
	err := s.db.QueryRowContext(ctx, query, a.Title, a.ArtistID, a.ReleaseDate).Scan(&a.ID)

	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	go s.setObjectToCache(ctx, strconv.Itoa(int(a.ID)), a)

	return a.ID, nil
}

// Add artist to database
func (s *store) CreateArtist(ctx context.Context, a *model.Artist) (int64, error) {
	const op = "repository.CatalogRepository.CreateArtist"

	query := "INSERT INTO artists(name) VALUES ($1) RETURNING id"
	err := s.db.QueryRowContext(ctx, query, a.Name).Scan(&a.ID)

	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	go s.setObjectToCache(ctx, strconv.Itoa(int(a.ID)), a)

	return a.ID, nil
}

// ArtistByID retrieves an artist by its ID from the database
func (s *store) ArtistByID(ctx context.Context, id int64) (*model.Artist, error) {
	const op = "repository.CatalogRepository.ArtistByID"

	key := strconv.Itoa(int(id))
	if cached, err := s.getObjectFromCache(ctx, key); err == nil && cached != nil {
		cachedArtist := cached.(*model.Artist)
		return cachedArtist, nil
	}

	query := "SELECT name FROM artists WHERE id = $1"

	row := s.db.QueryRowContext(ctx, query, id)

	var artist model.Artist

	err := row.Scan(&artist.Name)
	artist.ID = id

	if err != nil {
		return nil, s.handleError(op, err)
	}

	go s.setObjectToCache(ctx, key, &artist)
	return &artist, nil
}

// AlbumByID retrieves an album by its ID from the database
func (s *store) AlbumByID(ctx context.Context, id int64) (*model.Album, error) {
	const op = "repository.CatalogRepository.AlbumByID"

	key := strconv.Itoa(int(id))
	if cached, err := s.getObjectFromCache(ctx, key); err == nil && cached != nil {
		cachedAlbum := cached.(*model.Album)
		return cachedAlbum, nil
	}

	query := "SELECT title, artist_id, release_date FROM albums WHERE id = $1"

	row := s.db.QueryRowContext(ctx, query, id)

	var album model.Album

	err := row.Scan(&album.Title, &album.ArtistID, &album.ReleaseDate)
	album.ID = id

	if err != nil {
		return nil, s.handleError(op, err)
	}

	go s.setObjectToCache(ctx, key, &album)
	return &album, nil
}

// DeleteArtist deletes an artist from the database
func (s *store) DeleteArtist(ctx context.Context, id int64) (int64, error) {
	const op = "repository.ArtistRepository.DeleteArtist"

	query := "DELETE FROM artists WHERE id = $1"

	result, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return 0, s.handleError(op, err)
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return 0, s.handleError(op, err)
	}

	go s.deleteObjectFromCache(ctx, strconv.Itoa(int(id)))

	return id, nil
}

// DeleteAlbum deletes an album from the database
func (s *store) DeleteAlbum(ctx context.Context, id int64) (int64, error) {
	const op = "repository.AlbumRepository.DeleteAlbum"

	query := "DELETE FROM albums WHERE id = $1"

	result, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return 0, s.handleError(op, err)
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return 0, s.handleError(op, err)
	}

	go s.deleteObjectFromCache(ctx, strconv.Itoa(int(id)))

	return id, nil
}

func (s *store) handleError(op string, err error) error {
	if err == sql.ErrNoRows {
		return errors.NotFoundError(op, "track not found")
	}

	return errors.DatabaseError(op, err)
}
