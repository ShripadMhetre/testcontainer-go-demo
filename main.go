package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"testcontainers-demo/app"
	"testcontainers-demo/handlers"
	"testcontainers-demo/models"
	"testcontainers-demo/repositories"
)

func main() {
	router := gin.Default()

	// postgres connection
	db, err := repositories.ConnectToPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	// Perform database migration
	err = db.AutoMigrate(&models.Resource{})
	if err != nil {
		log.Fatal(err)
	}

	// repository test code
	// Create a user
	//newResource := &models.Resource{ID: uuid.New().String(), OfferId: "dunes", SiteGeoLocation: "us-east-1",
	//	AssetInfo: models.AssetInfo{
	//		AssetTag:    "bananas",
	//		AssetType:   "asdf",
	//		AssetFamily: "asdf",
	//		ServerType:  "adsf",
	//	}}
	//err = repositories.CreateResource(newResource)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Created Resource: %#v", newResource)

	// Query resource by ID
	//fetchedResource, err := repositories.GetResourceById(newResource.ResourceId)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Resource by ID: %#v", fetchedResource)

	// Routes
	router.GET("/metadata", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": app.Connections,
		})
	})

	router.GET("/resource/:resource-id", handlers.GetResourceHandler)
	router.POST("/resource", handlers.CreateResourceHandler)

	log.Fatal(router.Run(":8080"))
}
