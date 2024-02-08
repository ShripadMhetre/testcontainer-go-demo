package repositories

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

// Repository is the interface that wraps the basic operations with the Redis store.
type Repository struct {
	client *redis.Client
}

// NewRepository creates a new repository. It will receive a context and the Redis connection string.
func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	options, err := redis.ParseURL(connStr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to Redis: %v\n", err)
		return nil, err
	}

	cli := redis.NewClient(options)

	pong, err := cli.Ping(ctx).Result()
	if err != nil {
		// You probably want to retry here
		return nil, err
	}

	if pong != "PONG" {
		// You probably want to retry here
		return nil, err
	}

	return &Repository{client: cli}, nil
}
