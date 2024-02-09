package app

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

func CreatePostgreSQLContainer() (testcontainers.Container, error) {
	ctx := context.Background()
	c, err := postgres.RunContainer(ctx,
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
		return nil, err
	}

	postgresConn, err := c.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Postgres Connection String: ", postgresConn)

	Connections.PostgresURL = postgresConn
	return c, nil
}

func CreateRedisContainer() (testcontainers.Container, error) {
	ctx := context.Background()

	//c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	//	ContainerRequest: testcontainers.ContainerRequest{
	//		Image:        "redis:6-alpine",
	//		ExposedPorts: []string{"6379/tcp"},
	//		WaitingFor:   wait.ForLog("Ready to accept connections"),
	//	},
	//	Started: true,
	//})

	c, err := redis.RunContainer(ctx, testcontainers.WithImage("redis:6-alpine"))
	if err != nil {
		return nil, err
	}

	redisConn, err := c.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Redis Connection String: ", redisConn)

	Connections.RedisURL = redisConn
	return c, nil
}
