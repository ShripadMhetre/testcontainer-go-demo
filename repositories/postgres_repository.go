package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testcontainers-demo/app"
	"testcontainers-demo/models"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository() (*PostgresRepository, error) {
	//dsn := "user=username password=password dbname=dbname host=localhost port=5432 sslmode=disable"
	dsn := app.Connections.PostgresURL
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Perform database migration
	err = client.AutoMigrate(&models.Resource{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB Migration Done!!!")

	return &PostgresRepository{
		db: client,
	}, nil
}

func (r PostgresRepository) CreateResource(resource *models.Resource) error {
	result := r.db.Create(resource)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r PostgresRepository) GetResourceById(resourceId string) (*models.Resource, error) {
	var resource models.Resource
	result := r.db.First(&resource, "id = ?", resourceId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &resource, nil
}
