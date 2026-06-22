package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
}

// New creates and returns a new Redis store connection
func New(addr string) *RedisStore {
	return &RedisStore{
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: "",
			DB:       0,
			Protocol: 2,
		}),
	}
}

func (s *RedisStore) GetSuggestions(ctx context.Context, prefix string, limit int) ([]string, error) {
	// Query Redis for top suggestions matching the prefix
	suggestions, err := s.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   prefix,
		Start: 0,
		Stop:  limit - 1,
		Rev:   true,
	}).Result()
	if err != nil {
		return nil, err
	}

	return suggestions, nil
}

func (s *RedisStore) IncrementFrequency(ctx context.Context, query string) error {
	// Update frequency for all prefixes of the given query
	for i := range query {
		if err := s.rdb.ZIncrBy(ctx, query[:i+1], 1, query).Err(); err != nil {
			return err
		}
	}

	return nil
}
