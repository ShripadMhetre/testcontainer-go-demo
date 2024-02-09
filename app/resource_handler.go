package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testcontainers-demo/models"
	"testcontainers-demo/services"
)

type ResourceHandler struct {
	resourceService *services.ResourceService
}

func NewResourceHandler(resourceService *services.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

func (h ResourceHandler) CreateResourceHandler(c *gin.Context) {
	var resourceRequest models.CreateResourceRequest
	if err := c.ShouldBindJSON(&resourceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newResourceId, err := h.resourceService.CreateResource(c, resourceRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resourceId": newResourceId})
}

func (h ResourceHandler) GetResourceHandler(c *gin.Context) {
	resourceId := c.Param("resource-id")

	resource, err := h.resourceService.GetResource(c, resourceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resource)
}
