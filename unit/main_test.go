package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFilterComponentsByTypeAndBrand(t *testing.T) {
	components := getMockPCComponents()

	filtered := filterComponentsByTypeAndBrand(components, "CPU", "Intel")

	assert.Equal(t, 1, len(filtered))
	assert.Equal(t, "Intel", filtered[0].Brand)
	assert.Equal(t, "CPU", filtered[0].Type)
}

func TestSortComponentsByName(t *testing.T) {
	components := getMockPCComponents()

	sorted := sortComponents(components, "name")

	assert.Equal(t, "AMD Radeon RX 6800 XT", sorted[0].Name)
	assert.Equal(t, "Samsung 970 EVO Plus 1TB", sorted[len(sorted)-1].Name)
}

func TestPaginateComponents(t *testing.T) {
	components := getMockPCComponents()

	paginated := paginateComponents(components, "1", 3)

	assert.Equal(t, 3, len(paginated))
	assert.Equal(t, "Intel Core i9-10900K", paginated[0].Name)
	assert.Equal(t, "NVIDIA GeForce RTX 3080", paginated[1].Name)
	assert.Equal(t, "Corsair Vengeance RGB Pro 16GB", paginated[2].Name)
}

func TestFilteredComponentsHandler(t *testing.T) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	req, err := http.NewRequest("GET", "/filtered-components?type=CPU&brand=Intel&sort=name&page=1", nil)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.True(t, strings.Contains(resp.Body.String(), "Intel Core i9-10900K"))
}
