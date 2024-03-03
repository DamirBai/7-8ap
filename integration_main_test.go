package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFilteredComponentsIntegration(t *testing.T) {
	router := setupRouter()

	req, err := http.NewRequest("GET", "/filtered-components?type=CPU&brand=Intel&sort=name&page=1", nil)
	assert.NoError(t, err, "Error creating request")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Response status should be OK")
}

func setupRouter() http.Handler {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/components", func(c *gin.Context) {
		components := getMockPCComponents()

		c.HTML(http.StatusOK, "components.html", gin.H{
			"Components": components,
		})
	})

	router.GET("/filtered-components", func(c *gin.Context) {
		componentType := c.Query("type")
		brandFilter := c.Query("brand")
		sortBy := c.Query("sort")
		pageStr := c.Query("page")
		itemsPerPage := 3

		components := getMockPCComponents()

		filteredComponents := filterComponentsByTypeAndBrand(components, componentType, brandFilter)
		if len(filteredComponents) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No data available"})
			return
		}

		sortedComponents := sortComponents(filteredComponents, sortBy)

		paginatedComponents := paginateComponents(sortedComponents, pageStr, itemsPerPage)

		c.HTML(http.StatusOK, "components.html", gin.H{"Components": paginatedComponents})
	})

	return router
}
