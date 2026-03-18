package repository

import (
	"context"
	"encoding/json"
	"time"
	"track-service/internal/domain/model"

	"github.com/redis/go-redis/v9"
)

const (
	cacheTTL = 15 * time.Minute
)

func (s *store) getTrackFromCache(ctx context.Context, key string) (*model.Track, error) {
	data, err := s.db.GetRedis(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		s.db.Log().WarnContext(ctx, "failed to get from redis", "key", key, "error", err)
		return nil, nil
	}

	var track model.Track
	if err := json.Unmarshal(data, &track); err != nil {
		s.db.Log().WarnContext(ctx, "failed to unmarshal track from redis", "key", key, "error", err)
		return nil, nil
	}
	return &track, nil
}

func (s *store) setTrackToCache(ctx context.Context, key string, track *model.Track) {
	data, err := json.Marshal(track)
	if err != nil {
		return
	}
	if err := s.db.SetRedis(ctx, key, data, cacheTTL).Err(); err != nil {

	}
}

func (s *store) deleteUserFromCache(ctx context.Context, keys ...string) {
	for _, key := range keys {
		if err := s.db.DelRedis(ctx, key).Err(); err != nil {
			s.db.Log().WarnContext(ctx, "failed to delete redis key", "key", key, "error", err)
		}
	}
}
