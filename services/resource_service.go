package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"testcontainers-demo/models"
	"testcontainers-demo/repositories"
)

type ResourceService struct {
	pgRepository    *repositories.PostgresRepository
	redisRepository *repositories.RedisRepository
}

func NewResourceService(pgRepo *repositories.PostgresRepository, redisRepo *repositories.RedisRepository) *ResourceService {
	return &ResourceService{
		pgRepository:    pgRepo,
		redisRepository: redisRepo,
	}
}

func (h ResourceService) CreateResource(c context.Context, resourceRequest models.CreateResourceRequest) (string, error) {
	newResourceId := uuid.New().String()
	newResource := &models.Resource{ID: newResourceId, OfferId: "dunes", SiteGeoLocation: resourceRequest.SiteGeoLocation,
		AssetInfo: resourceRequest.AssetInfo}
	err := h.pgRepository.CreateResource(newResource)
	if err != nil {
		return "", err
	}

	// Cache the resource in Redis
	err = h.redisRepository.CacheResource(c, newResourceId, *newResource)
	if err != nil {
		return "", err
	}

	return newResourceId, nil
}

func (h ResourceService) GetResource(c context.Context, resourceId string) (*models.Resource, error) {
	// Check if the resource is cached in Redis
	cachedResource, err := h.redisRepository.GetResource(c, resourceId)
	if err == nil {
		fmt.Println("Resource found in Redis for resourceId:", resourceId)
		return cachedResource, nil
	}

	// If the resource is not cached in Redis, fetch it from the database
	resource, err := h.pgRepository.GetResourceById(resourceId)
	if err != nil {
		return nil, err
	}

	fmt.Println("Resource fetched from DB:", resourceId)
	return resource, nil
}
