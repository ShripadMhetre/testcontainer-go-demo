package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"testcontainers-demo/models"
	"testcontainers-demo/repositories"
)

func CreateResourceHandler(c *gin.Context) {
	var resourceRequest models.CreateResourceRequest
	if err := c.ShouldBindJSON(&resourceRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newResourceId := uuid.New().String()
	newResource := &models.Resource{ID: newResourceId, OfferId: "dunes", SiteGeoLocation: resourceRequest.SiteGeoLocation,
		AssetInfo: resourceRequest.AssetInfo}
	err := repositories.CreateResource(newResource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resourceId": newResourceId})
}

func GetResourceHandler(c *gin.Context) {
	resourceId := c.Param("resource-id")

	resource, err := repositories.GetResourceById(resourceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resource)
}
