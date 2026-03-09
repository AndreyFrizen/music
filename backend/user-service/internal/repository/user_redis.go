package repository

import (
	"context"
	"encoding/json"
	"time"
	modeluser "user-service/internal/domain/model"

	"github.com/redis/go-redis/v9"
)

const (
	cacheTTL = 15 * time.Minute
)

func (s *store) getUserFromCache(ctx context.Context, key string) (*modeluser.User, error) {
	data, err := s.db.GetRedis(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		s.db.Log().WarnContext(ctx, "failed to get from redis", "key", key, "error", err)
		return nil, nil
	}

	var user modeluser.User
	if err := json.Unmarshal(data, &user); err != nil {
		s.db.Log().WarnContext(ctx, "failed to unmarshal user from redis", "key", key, "error", err)
		return nil, nil
	}
	return &user, nil
}

func (s *store) setUserToCache(ctx context.Context, key string, user *modeluser.User) {
	data, err := json.Marshal(user)
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
