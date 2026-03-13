package services

import (
	"context"
	"log/slog"
	"track-service/internal/domain/errors"
	"track-service/internal/domain/model"

	"github.com/go-playground/validator/v10"
)

type service struct {
	repo     TrackRepository
	log      *slog.Logger
	validate *validator.Validate
}

func NewService(repo TrackRepository, log *slog.Logger) *service {
	return &service{
		repo:     repo,
		log:      log,
		validate: validator.New(),
	}
}

type TrackRepository interface {
	CreateTrack(ctx context.Context, t *model.NewTrack) (int64, error)
	TrackByID(ctx context.Context, id int64) (*model.Track, error)
	UpdateTrack(ctx context.Context, t *model.Track) error
	DeleteTrack(ctx context.Context, id int64) error
}

// CreateTrack creates a track in the database
func (s *service) CreateTrack(ctx context.Context, req *CreateTrackRequest) (*CreateTrackResponse, error) {
	const op = "service.UserService.Register"

	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(op, map[string]string{
			"track": "invalid track",
		})
	}

	t := &model.NewTrack{
		Title:    req.Title,
		Duration: req.Duration,
		ArtistID: req.ArtistID,
		AlbumID:  req.AlbumID,
	}

	track, err := s.repo.CreateTrack(ctx, t)
	if err != nil {
		s.log.Error(op, "failed to create track", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track created successfully", req.Title)

	return &CreateTrackResponse{ID: track}, nil
}

// TrackByID retrieves a track by its ID
func (s *service) TrackByID(ctx context.Context, req *GetTrackRequest) (*GetTrackResponse, error) {
	const op = "service.TrackByID"

	track, err := s.repo.TrackByID(ctx, req.ID)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "track not found")
		}
		s.log.Error(op, "failed to get track", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track retrieved successfully", track.ID)

	return &GetTrackResponse{
		ID:       track.ID,
		Title:    track.Title,
		Duration: track.Duration,
		AlbumID:  track.AlbumID,
		ArtistID: track.ArtistID,
		AudioURL: track.AudioURL,
	}, nil
}

// UpdateTrack updates a track by its ID
func (s *service) UpdateTrack(ctx context.Context, req *UpdateTrackRequest) (*UpdateTrackResponse, error) {
	const op = "service.UpdateTrack"

	t := &model.Track{
		ID:       req.ID,
		Title:    req.Title,
		Duration: req.Duration,
		ArtistID: req.ArtistID,
		AlbumID:  req.AlbumID,
	}

	err := s.repo.UpdateTrack(ctx, t)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "track not found")
		}
		s.log.Error(op, "failed to update track", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track updated successfully", req.ID)

	return &UpdateTrackResponse{ID: req.ID}, nil
}

// DeleteTrack deletes a track by its ID
func (s *service) DeleteTrack(ctx context.Context, req *DeleteTrackRequest) (*DeleteTrackResponse, error) {
	const op = "service.DeleteTrack"

	err := s.repo.DeleteTrack(ctx, req.ID)
	if err != nil {
		s.log.Error(op, "failed to delete track", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track deleted successfully", req.ID)

	return &DeleteTrackResponse{Success: true}, nil
}
