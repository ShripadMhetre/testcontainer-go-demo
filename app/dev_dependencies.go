//go:build dev || e2e
// +build dev e2e

package app

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	startupDependenciesFns := []func() (testcontainers.Container, error){
		CreatePostgreSQLContainer,
		CreateRedisContainer,
	}

	runtimeDependencies := make([]testcontainers.Container, 0, len(startupDependenciesFns))

	for _, fn := range startupDependenciesFns {
		c, err := fn()
		if err != nil {
			panic(err)
		}
		runtimeDependencies = append(runtimeDependencies, c)
	}

	// register a graceful shutdown to stop the dependencies when the application is stopped
	// only in development mode
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		// also use the shutdown function when the SIGTERM or SIGINT signals are received
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v\n", sig)
		err := shutdownDependencies(runtimeDependencies...)
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()
}

// helper function to stop the dependencies
func shutdownDependencies(containers ...testcontainers.Container) error {
	ctx := context.Background()
	for _, c := range containers {
		err := c.Terminate(ctx)
		if err != nil {
			return fmt.Errorf("failed to terminate container: %w", err)
		}
	}

	return nil
}

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
