package store

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr string) (*RedisStore, error) {
	opts, err := redis.ParseURL(addr)
	if err != nil {
		// Treat addr as host:port if not a full URL
		opts = &redis.Options{Addr: addr}
	}

	client := redis.NewClient(opts)

	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisStore{client: client, ctx: ctx}, nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

func (s *RedisStore) Set(alias, originalURL string, expiresIn int) error {
	key := "url:" + alias
	clicksKey := "clicks:" + alias

	if expiresIn > 0 {
		err := s.client.Set(s.ctx, key, originalURL, time.Duration(expiresIn)*time.Second).Err()
		if err != nil {
			return err
		}
		s.client.Expire(s.ctx, clicksKey, time.Duration(expiresIn)*time.Second)
	} else {
		err := s.client.Set(s.ctx, key, originalURL, 0).Err()
		if err != nil {
			return err
		}
	}

	s.client.Set(s.ctx, clicksKey, 0, 0)
	return nil
}

func (s *RedisStore) Get(alias string) (string, error) {
	return s.client.Get(s.ctx, "url:"+alias).Result()
}

func (s *RedisStore) Exists(alias string) bool {
	val, err := s.client.Exists(s.ctx, "url:"+alias).Result()
	return err == nil && val > 0
}

func (s *RedisStore) IncrementClick(alias string) {
	s.client.Incr(s.ctx, "clicks:"+alias)
}