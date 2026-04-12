package repository

import (
	"collection-service/proto/models"
	"context"
	"encoding/json"

	"time"

	"github.com/redis/go-redis/v9"
)

const (
	cacheTTL = 15 * time.Minute
)

func (s *store) getAlbumsFromCache(ctx context.Context, key string) ([]*models.Album, error) {
	var albums []*models.Album
	data, err := s.db.GetRedis(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		s.db.Log().WarnContext(ctx, "failed to get from redis", "key", key, "error", err)
		return nil, nil
	}
	if err := json.Unmarshal(data, &albums); err != nil {
		s.db.Log().WarnContext(ctx, "failed to unmarshal album from redis", "key", key, "error", err)
		return nil, nil
	}
	return albums, nil
}

func (s *store) setAlbumToCache(ctx context.Context, key string, album *models.Album) {
	data, err := json.Marshal(album)
	if err != nil {
		return
	}
	if err := s.db.SetRedis(ctx, key, data, cacheTTL).Err(); err != nil {
		return
	}
}

func (s *store) deleteAlbumFromCache(ctx context.Context, key string) {
	if err := s.db.DelRedis(ctx, key).Err(); err != nil {
		s.db.Log().WarnContext(ctx, "failed to delete redis key", "key", key, "error", err)
	}
}
