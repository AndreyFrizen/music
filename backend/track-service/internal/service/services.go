package services

import (
	"context"
	"log/slog"
	"track-service/internal/domain/errors"

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
	CreateTrack(ctx context.Context, req *CreateTrackRequest) (*CreateTrackResponse, error)
	TrackByID(ctx context.Context, req *GetTrackRequest) (*GetTrackResponse, error)
	UpdateTrack(ctx context.Context, req *UpdateTrackRequest) (*UpdateTrackResponse, error)
	DeleteTrack(ctx context.Context, req *DeleteTrackRequest) (*DeleteTrackResponse, error)
}

// CreateTrack creates a track in the database
func (s *service) CreateTrack(ctx context.Context, req *CreateTrackRequest) (*CreateTrackResponse, error) {
	const op = "service.UserService.Register"

	if err := s.validate.Struct(req); err != nil {
		return nil, errors.ValidationError(op, map[string]string{
			"track": "invalid track",
		})
	}

	track, err := s.repo.CreateTrack(ctx, req)
	if err != nil {
		s.log.Error(op, "failed to create track", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track created successfully", req.Title)

	return &CreateTrackResponse{ID: track.ID}, nil
}

// TrackByID retrieves a track by its ID
func (s *service) TrackByID(ctx context.Context, req *GetTrackRequest) (*GetTrackResponse, error) {
	const op = "service.TrackByID"

	track, err := s.repo.TrackByID(ctx, req)

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "track not found")
		}
		s.log.Error(op, "failed to get track", "track_id", req.ID, "error", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track retrieved successfully",
		"op", op,
		"track_id", track.ID,
	)

	return track, nil
}

// UpdateTrack updates a track by its ID
func (s *service) UpdateTrack(ctx context.Context, req *UpdateTrackRequest) (*UpdateTrackResponse, error) {
	const op = "service.UpdateTrack"

	track, err := s.repo.UpdateTrack(ctx, req)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "track not found")
		}
		s.log.Error(op, "failed to update track", "track_id", req.ID, "error", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track updated successfully",
		"op", op,
		"track_id", track.ID,
	)

	return track, nil
}

func (s *service) DeleteTrack(ctx context.Context, req *DeleteTrackRequest) (*DeleteTrackResponse, error) {
	const op = "service.DeleteTrack"

	sc, err := s.repo.DeleteTrack(ctx, req)
	if err != nil {
		s.log.Error(op, "failed to delete track", "track_id", req.ID, "error", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track deleted successfully",
		"op", op,
		"track_id", req.ID,
	)

	return sc, nil
}
