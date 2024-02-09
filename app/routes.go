package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"testcontainers-demo/repositories"
	"testcontainers-demo/services"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Postgres repository instantiation
	pgRepository, err := repositories.NewPostgresRepository(Connections.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}

	// Redis repository instantiation
	redisRepository, err := repositories.NewRedisRepository(context.Background(), Connections.RedisURL)
	if err != nil {
		log.Fatal(err)
	}

	// Resource service & handler instantiation
	resourceService := services.NewResourceService(pgRepository, redisRepository)
	resourceHandler := NewResourceHandler(resourceService)

	// Helper route
	router.GET("/metadata", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"connections": Connections,
		})
	})

	// Business routes
	router.GET("/resources/:resource-id", resourceHandler.GetResourceHandler)
	router.POST("/resources", resourceHandler.CreateResourceHandler)

	return router
}
