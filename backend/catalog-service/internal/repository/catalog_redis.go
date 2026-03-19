package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type ObjInterface interface {
	GetID() int64
}

const (
	cacheTTL = 15 * time.Minute
)

func (s *store) getObjectFromCache(ctx context.Context, key string) (ObjInterface, error) {
	data, err := s.db.GetRedis(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		s.db.Log().WarnContext(ctx, "failed to get from redis", "key", key, "error", err)
		return nil, nil
	}

	var obj ObjInterface
	if err := json.Unmarshal(data, &obj); err != nil {
		s.db.Log().WarnContext(ctx, "failed to unmarshal object from redis", "key", key, "error", err)
		return nil, nil
	}
	return obj, nil
}

func (s *store) setObjectToCache(ctx context.Context, key string, obj ObjInterface) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	if err := s.db.SetRedis(ctx, key, data, cacheTTL).Err(); err != nil {

	}
}

func (s *store) deleteObjectFromCache(ctx context.Context, keys ...string) {
	for _, key := range keys {
		if err := s.db.DelRedis(ctx, key).Err(); err != nil {
			s.db.Log().WarnContext(ctx, "failed to delete redis key", "key", key, "error", err)
		}
	}
}
