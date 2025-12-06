package redis

import (
	"chat2pay/config/yaml"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (*string, error)
	Set(ctx context.Context, key string, value string) (*string, error)
}

type redisClient struct {
	client *redis.Client
}

func NewRedis(cfg *yaml.Config) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%s`, cfg.Redis.Host, cfg.Redis.Port), // Redis server address
		Password: cfg.Redis.Password,                                   // No password set
		DB:       cfg.Redis.DB,
	},
	)

	return &redisClient{client: client}
}

func (r *redisClient) Get(ctx context.Context, key string) (*string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}

		return nil, err
	}

	return &val, nil
}

func (r *redisClient) Set(ctx context.Context, key string, value string) (*string, error) {
	val, err := r.client.Set(ctx, key, value, 2*time.Hour).Result()

	if err != nil {
		return nil, err
	}

	return &val, nil
}
