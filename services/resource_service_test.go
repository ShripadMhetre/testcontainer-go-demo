package services

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testcontainers-demo/models"
	"testcontainers-demo/repositories"
	"testing"
	"time"
)

type ResourceServiceTestSuite struct {
	suite.Suite
	ctx             context.Context
	pgContainer     testcontainers.Container
	redisContainer  testcontainers.Container
	pgRepository    *repositories.PostgresRepository
	redisRepository *repositories.RedisRepository
	resourceService *ResourceService

	resourceId string
}

func (suite *ResourceServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Postgres container
	pgContainer, err := postgres.RunContainer(suite.ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		//postgres.WithInitScripts(filepath.Join("./app", "db-data", "dev-db.sql")),
		postgres.WithDatabase("resource-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		log.Fatal(err)
	}
	postgresConn, err := pgContainer.ConnectionString(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContainer = pgContainer
	pgRepository, err := repositories.NewPostgresRepository(postgresConn)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgRepository = pgRepository

	// Redis container
	redisContainer, err := redis.RunContainer(suite.ctx, testcontainers.WithImage("redis:6-alpine"))
	if err != nil {
		log.Fatal(err)
	}

	redisConn, err := redisContainer.ConnectionString(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.redisContainer = redisContainer
	redisRepository, err := repositories.NewRedisRepository(suite.ctx, redisConn)
	if err != nil {
		log.Fatal(err)
	}
	suite.redisRepository = redisRepository

	// Resource service instantiation
	resourceService := NewResourceService(pgRepository, redisRepository)
	suite.resourceService = resourceService
}

func (suite *ResourceServiceTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("Error terminating postgres container: %s", err)
	}

	if err := suite.redisContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("Error terminating redis container: %s", err)
	}
}

func TestResourceTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceServiceTestSuite))
}

func (suite *ResourceServiceTestSuite) TestCreateResource() {
	t := suite.T()

	resourceRequest := models.CreateResourceRequest{
		SiteGeoLocation: "test-location",
		AssetInfo: models.AssetInfo{
			AssetTag:    "assetTag",
			AssetType:   "assetType",
			AssetFamily: "assetFamily",
			ServerType:  "serverType",
		},
	}

	resourceId, err := suite.resourceService.CreateResource(suite.ctx, resourceRequest)
	suite.resourceId = resourceId

	assert.NoError(t, err)
	assert.Equal(t, 36, len(resourceId))
}

func (suite *ResourceServiceTestSuite) TestGetResource() {
	t := suite.T()

	resource, err := suite.resourceService.GetResource(suite.ctx, suite.resourceId)
	assert.NoError(t, err)
	assert.NotNil(t, resource)
}
