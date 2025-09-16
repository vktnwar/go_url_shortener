package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	Client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{Client: client}
}

func (r *RedisRepository) SetURL(ctx context.Context, shortID, originalURL string, ttl time.Duration) error {
	return r.Client.Set(ctx, shortID, originalURL, ttl).Err()
}

func (r *RedisRepository) GetURL(ctx context.Context, shortID string) (string, error) {
	return r.Client.Get(ctx, shortID).Result()
}
