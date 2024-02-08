package repositories

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testcontainers-demo/app"
	"testcontainers-demo/models"
)

var db *gorm.DB

func ConnectToPostgreSQL() (*gorm.DB, error) {
	//dsn := "user=username password=password dbname=dbname host=localhost port=5432 sslmode=disable"
	dsn := app.Connections.PostgresURL
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db = client
	return db, nil
}

func CreateResource(resource *models.Resource) error {
	result := db.Create(resource)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetResourceById(resourceId string) (*models.Resource, error) {
	var resource models.Resource
	result := db.First(&resource, "id = ?", resourceId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &resource, nil
}
