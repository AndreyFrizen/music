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
func (s *service) CreateTrack(ctx context.Context, req *model.CreateTrackRequest) (*model.CreateTrackResponse, error) {
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

	return &model.CreateTrackResponse{ID: track}, nil
}

// TrackByID retrieves a track by its ID
func (s *service) TrackByID(ctx context.Context, req *model.GetTrackRequest) (*model.GetTrackResponse, error) {
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

	return &model.GetTrackResponse{
		Track: track,
	}, nil
}

// UpdateTrack updates a track by its ID
func (s *service) UpdateTrack(ctx context.Context, req *model.UpdateTrackRequest) (*model.UpdateTrackResponse, error) {
	const op = "service.UpdateTrack"

	t := req.Track

	err := s.repo.UpdateTrack(ctx, t)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundError(op, "track not found")
		}
		s.log.Error(op, "failed to update track", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track updated successfully", t.ID)

	return &model.UpdateTrackResponse{ID: t.ID}, nil
}

// DeleteTrack deletes a track by its ID
func (s *service) DeleteTrack(ctx context.Context, req *model.DeleteTrackRequest) (*model.DeleteTrackResponse, error) {
	const op = "service.DeleteTrack"

	err := s.repo.DeleteTrack(ctx, req.ID)
	if err != nil {
		s.log.Error(op, "failed to delete track", err)
		return nil, errors.InternalError(op, err)
	}

	s.log.Info(op, "track deleted successfully", req.ID)

	return &model.DeleteTrackResponse{Success: true}, nil
}
