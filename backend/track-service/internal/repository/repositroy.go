package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"track-service/internal/app/database"
	"track-service/internal/domain/errors"
	"track-service/internal/domain/model"
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
func (s *store) CreateTrack(ctx context.Context, t *model.NewTrack) (int64, error) {
	const op = "repository.TrackRepository.CreateTrack"

	var id int64
	err := s.db.QueryRowContext(ctx,
		"INSERT INTO tracks(title, duration, artist_id, album_id) VALUES ($1, $2, $3, $4) RETURNING id",
		t.Title, t.Duration, t.ArtistID, t.AlbumID).Scan(&id)

	url := fmt.Sprintf("%s/%d", "tracks", id)

	_, err = s.db.ExecContext(ctx, "UPDATE tracks SET audio_url = $1 WHERE id = $2", url, id)

	if err != nil {
		return 0, errors.DatabaseError(op, err)
	}

	go func() {
		t, err := s.TrackByID(ctx, id)
		if err == nil {
			s.setTrackToCache(ctx, strconv.Itoa(int(id)), t)
		}
	}()

	return id, nil
}

// TrackByID retrieves a track by its ID from the database
func (s *store) TrackByID(ctx context.Context, id int64) (*model.Track, error) {
	const op = "repository.TrackRepository.TrackByID"

	key := strconv.Itoa(int(id))
	if cached, err := s.getTrackFromCache(ctx, key); err == nil && cached != nil {
		return cached, nil
	}

	query := "SELECT id, title, duration, audio_url, artist_id FROM tracks WHERE id = $1"

	row := s.db.QueryRowContext(ctx, query, id)

	var track model.Track

	err := row.Scan(&track.ID, &track.Title, &track.Duration, &track.AudioURL, &track.ArtistID, &track.AlbumID)

	if err != nil {
		return nil, s.handleError(op, err)
	}

	s.setTrackToCache(ctx, key, &track)
	return &track, nil
}

// UpdateTrack updates a track in the database
func (s *store) UpdateTrack(ctx context.Context, t *model.Track) error {
	const op = "repository.TrackRepository.UpdateTrack"

	result, err := s.db.ExecContext(ctx,
		"UPDATE tracks SET title = $1, duration = $2, artist_id = $3 WHERE id = $4",
		t.Title, t.Duration, t.ArtistID, t.ID)

	if err != nil {
		return s.handleError(op, err)
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return s.handleError(op, err)
	}

	go s.setTrackToCache(ctx, strconv.Itoa(int(t.ID)), t)

	return nil
}

// DeleteTrack deletes a track from the database
func (s *store) DeleteTrack(ctx context.Context, id int64) error {
	const op = "repository.TrackRepository.DeleteTrack"

	query := "DELETE FROM tracks WHERE id = $1"

	result, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return s.handleError(op, err)
	}

	rows := result.RowsAffected()

	if rows == 0 {
		return s.handleError(op, err)
	}

	return nil
}

func (s *store) handleError(op string, err error) error {
	if err == sql.ErrNoRows {
		return errors.NotFoundError(op, "track not found")
	}

	return errors.DatabaseError(op, err)
}
