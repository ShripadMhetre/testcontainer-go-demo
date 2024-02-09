package handlers

//
//import (
//	"context"
//	"github.com/stretchr/testify/suite"
//	"log"
//	"testcontainers-demo/app"
//	"testcontainers-demo/repositories"
//	"testing"
//)
//
//type ResourceTestSuite struct {
//	suite.Suite
//	ctx             context.Context
//	pgRepository    *repositories.PostgresRepository
//	redisRepository *repositories.RedisRepository
//}
//
//func (suite *ResourceTestSuite) SetupSuite() {
//	suite.ctx = context.Background()
//
//	// Create a new Postgres container and repository
//	_, err := app.CreatePostgreSQLContainer()
//	if err != nil {
//		log.Fatal(err)
//	}
//	pgRepository, err := repositories.NewPostgresRepository()
//	if err != nil {
//		log.Fatal(err)
//	}
//	suite.pgRepository = pgRepository
//
//	// Create a new Redis container and repository
//	_, err = app.CreateRedisContainer()
//	if err != nil {
//		log.Fatal(err)
//	}
//	redisRepository, err := repositories.NewRedisRepository(suite.ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	suite.redisRepository = redisRepository
//}
//
//func (suite *ResourceTestSuite) TearDownSuite() {
//	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
//		log.Fatalf("error terminating postgres container: %s", err)
//	}
//
//	//defer func() {
//	//	if err := redisContainer.Terminate(ctx); err != nil {
//	//		log.Fatalf("error terminating redis container: %s", err)
//	//	}
//	//}()
//}
//
//func TestResourceTestSuite(t *testing.T) {
//	suite.Run(t, new(ResourceTestSuite))
//}
