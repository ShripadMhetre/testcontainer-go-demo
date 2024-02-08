package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"testcontainers-demo/app"
	"testcontainers-demo/models"
)

// RedisRepository is the interface that wraps the basic operations with the Redis store.
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new repository. It will receive a context and the Redis connection string.
func NewRedisRepository(ctx context.Context) (*RedisRepository, error) {
	fmt.Println("REDIS URL: ", app.Connections.RedisURL)
	options, err := redis.ParseURL(app.Connections.RedisURL)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to Redis: %v\n", err)
		return nil, err
	}

	cli := redis.NewClient(options)

	pong, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	if pong != "PONG" {
		return nil, err
	}

	return &RedisRepository{client: cli}, nil
}

// CacheResource caches a resource in the Redis store.
func (r RedisRepository) CacheResource(ctx context.Context, key string, value models.Resource) error {
	// Marshal the resource to JSON
	resourceBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, resourceBytes, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetResource retrieves a resource from the Redis store.
func (r RedisRepository) GetResource(ctx context.Context, key string) (*models.Resource, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cachedResource models.Resource
	err = json.Unmarshal([]byte(val), &cachedResource)
	if err != nil {
		return nil, err
	}

	return &cachedResource, nil
}
