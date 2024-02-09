package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"testcontainers-demo/app"
	"testcontainers-demo/models"
	"testing"
)

func TestRedisRepository(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := redis.RunContainer(ctx, testcontainers.WithImage("docker.io/redis:6-alpine"))
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate redis container: %s", err)
		}
	})

	connStr, err := redisContainer.ConnectionString(ctx)
	require.NoError(t, err)

	app.Connections.RedisURL = connStr
	repo, err := NewRedisRepository(ctx)
	require.NoError(t, err)
	assert.NotNil(t, repo)

	testResourceId := uuid.New().String()
	t.Run("Cache Resource", func(t *testing.T) {
		newResource := &models.Resource{ID: testResourceId, OfferId: "dunes", SiteGeoLocation: "us-east-1",
			AssetInfo: models.AssetInfo{
				AssetTag:    "assetTag",
				AssetType:   "assetType",
				AssetFamily: "assetFamily",
				ServerType:  "serverType",
			}}

		err := repo.CacheResource(ctx, testResourceId, *newResource)
		assert.NoError(t, err)
	})

	t.Run("Fetch Cached Resource", func(t *testing.T) {
		cachedResource, err := repo.GetResource(ctx, testResourceId)
		assert.NoError(t, err)
		assert.NotNil(t, cachedResource)
		assert.Equal(t, testResourceId, cachedResource.ID)

	})
}
