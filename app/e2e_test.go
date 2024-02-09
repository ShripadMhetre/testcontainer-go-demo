//go:build e2e
// +build e2e

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testcontainers-demo/models"
	"testing"
)

func TestMetadataRouteWithDependencies(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metadata", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// the "GET /" endpoint returns a JSON with metadata including the connection strings for the dependencies
	var response struct {
		Connections Metadata `json:"connections"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// assert that the different connection strings are set
	assert.True(t, strings.Contains(response.Connections.RedisURL, "redis://localhost:"),
		fmt.Sprintf("expected %s to be a Redis URL", response.Connections.RedisURL))
	assert.True(t, strings.Contains(response.Connections.PostgresURL, "postgres://postgres:postgres@localhost:"),
		fmt.Sprintf("expected %s to be a Postgres URL", response.Connections.PostgresURL))
}

func TestRoutesWithDependencies(t *testing.T) {
	router := SetupRouter()

	var testResourceId string
	t.Run("POST /resources", func(t *testing.T) {
		resourceReq := models.CreateResourceRequest{
			SiteGeoLocation: "test-location",
			AssetInfo: models.AssetInfo{
				AssetTag:    "assetTag",
				AssetType:   "assetType",
				AssetFamily: "assetFamily",
				ServerType:  "serverType",
			},
		}

		body, err := json.Marshal(resourceReq)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/resources", bytes.NewReader(body))
		require.NoError(t, err)

		// we need to set the content type header because we are sending a body
		req.Header.Add("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// get the resourceId from the response
		var response struct {
			ResourceId string `json:"resourceId"`
		}
		err = json.Unmarshal(w.Body.Bytes(), &response)

		testResourceId = response.ResourceId

		// we are receiving a 200 because the ratings repository is started
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET /resource/:resource-id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/resources/"+testResourceId, nil)
		require.NoError(t, err)
		router.ServeHTTP(w, req)

		// we are receiving a 200 because the ratings repository is started
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
