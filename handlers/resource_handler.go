package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"testcontainers-demo/models"
	"testcontainers-demo/repositories"
)

type ResourceHandler struct {
	resourceRepository *repositories.PostgresRepository
}

func NewResourceHandler(resourceRepo *repositories.PostgresRepository) *ResourceHandler {
	return &ResourceHandler{
		resourceRepository: resourceRepo,
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
	err := h.resourceRepository.CreateResource(newResource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resourceId": newResourceId})
}

func (h ResourceHandler) GetResourceHandler(c *gin.Context) {
	resourceId := c.Param("resource-id")

	resource, err := h.resourceRepository.GetResourceById(resourceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resource)
}
