package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"testcontainers-demo/app"
	"testcontainers-demo/handlers"
	"testcontainers-demo/repositories"
)

func main() {
	router := gin.Default()

	// Postgres repository instantiation
	pgRepository, err := repositories.NewPostgresRepository()
	if err != nil {
		log.Fatal(err)
	}

	// Helper route
	router.GET("/metadata", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": app.Connections,
		})
	})

	resourceHandler := handlers.NewResourceHandler(pgRepository)

	// Business routes
	router.GET("/resource/:resource-id", resourceHandler.GetResourceHandler)
	router.POST("/resource", resourceHandler.CreateResourceHandler)

	log.Fatal(router.Run(":8080"))
}
