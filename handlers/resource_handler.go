package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"testcontainers-demo/models"
	"testcontainers-demo/repositories"
)

type ResourceHandler struct {
	pgRepository    *repositories.PostgresRepository
	redisRepository *repositories.RedisRepository
}

func NewResourceHandler(pgRepo *repositories.PostgresRepository, redisRepo *repositories.RedisRepository) *ResourceHandler {
	return &ResourceHandler{
		pgRepository:    pgRepo,
		redisRepository: redisRepo,
	}
}

func (h ResourceHandler) CreateResourceHandler(c *gin.Context) {
	var resourceRequest models.CreateResourceRequest
	if err := c.ShouldBindJSON(&resourceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newResourceId := uuid.New().String()
	newResource := &models.Resource{ID: newResourceId, OfferId: "dunes", SiteGeoLocation: resourceRequest.SiteGeoLocation,
		AssetInfo: resourceRequest.AssetInfo}
	err := h.pgRepository.CreateResource(newResource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cache the resource in Redis
	err = h.redisRepository.CacheResource(c, newResourceId, *newResource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resourceId": newResourceId})
}

func (h ResourceHandler) GetResourceHandler(c *gin.Context) {
	resourceId := c.Param("resource-id")

	// Check if the resource is cached in Redis
	cachedResource, err := h.redisRepository.GetResource(c, resourceId)
	if err == nil {
		fmt.Println("Resource found in Redis for resourceId:", resourceId)
		c.JSON(http.StatusOK, cachedResource)
		return
	}

	// If the resource is not cached in Redis, fetch it from the database
	resource, err := h.pgRepository.GetResourceById(resourceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Resource fetched from DB:", resourceId)
	c.JSON(http.StatusOK, resource)
}
