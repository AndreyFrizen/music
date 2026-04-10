package repository

import (
	"collection-service/internal/domain/model"
	"context"
	"encoding/json"

	"time"

	"github.com/redis/go-redis/v9"
)

const (
	cacheTTL = 15 * time.Minute
)

func (s *store) getPlaylistFromCache(ctx context.Context, key string) (*model.Playlist, error) {
	data, err := s.db.GetRedis(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		s.db.Log().WarnContext(ctx, "failed to get from redis", "key", key, "error", err)
		return nil, nil
	}

	var playlist model.Playlist
	if err := json.Unmarshal(data, &playlist); err != nil {
		s.db.Log().WarnContext(ctx, "failed to unmarshal playlist from redis", "key", key, "error", err)
		return nil, nil
	}
	return &playlist, nil
}

func (s *store) setPlaylistToCache(ctx context.Context, key string, playlist *model.Playlist) {
	data, err := json.Marshal(playlist)
	if err != nil {
		return
	}
	if err := s.db.SetRedis(ctx, key, data, cacheTTL).Err(); err != nil {

	}
}

func (s *store) deletePlaylistFromCache(ctx context.Context, keys ...string) {
	for _, key := range keys {
		if err := s.db.DelRedis(ctx, key).Err(); err != nil {
			s.db.Log().WarnContext(ctx, "failed to delete redis key", "key", key, "error", err)
		}
	}
}
