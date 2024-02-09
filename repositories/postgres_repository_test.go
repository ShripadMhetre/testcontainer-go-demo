package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testcontainers-demo/models"
	"testing"
	"time"
)

func TestPostgresRepository(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		//postgres.WithInitScripts(filepath.Join("./app", "db-data", "dev-db.sql")),
		postgres.WithDatabase("resource-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate postgreSQL container: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx)
	require.NoError(t, err)

	repo, err := NewPostgresRepository(connStr)
	require.NoError(t, err)
	assert.NotNil(t, repo)

	testResourceId := uuid.New().String()
	t.Run("Create Resource in DB", func(t *testing.T) {
		newResource := &models.Resource{ID: testResourceId, OfferId: "dunes", SiteGeoLocation: "us-east-1",
			AssetInfo: models.AssetInfo{
				AssetTag:    "assetTag",
				AssetType:   "assetType",
				AssetFamily: "assetFamily",
				ServerType:  "serverType",
			}}

		err := repo.CreateResource(newResource)
		assert.NoError(t, err)
	})

	t.Run("Fetch Resource from DB", func(t *testing.T) {
		cachedResource, err := repo.GetResourceById(testResourceId)
		assert.NoError(t, err)
		assert.NotNil(t, cachedResource)
		assert.Equal(t, testResourceId, cachedResource.ID)
	})
}
