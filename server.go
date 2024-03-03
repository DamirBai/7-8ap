package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 3)

type PCComponent struct {
	Type     string
	Brand    string
	Name     string
	ImageURL string
	Price    int16
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func filterComponentsByTypeAndBrand(components []PCComponent, componentType, brandFilter string) []PCComponent {
	if componentType == "" && brandFilter == "" {
		return components
	}

	var filteredComponents []PCComponent
	for _, comp := range components {
		if (componentType == "" || comp.Type == componentType) && (brandFilter == "" || comp.Brand == brandFilter) {
			filteredComponents = append(filteredComponents, comp)
		}
	}
	return filteredComponents
}

func sortComponents(components []PCComponent, sortBy string) []PCComponent {
	switch sortBy {
	case "name":
		sort.Slice(components, func(i, j int) bool {
			return components[i].Name < components[j].Name
		})
	case "price":
		sort.Slice(components, func(i, j int) bool {
			return components[i].Price < components[j].Price
		})
	}
	return components
}

func paginateComponents(components []PCComponent, pageStr string, itemsPerPage int) []PCComponent {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	startIdx := (page - 1) * itemsPerPage
	endIdx := startIdx + itemsPerPage

	if startIdx >= len(components) {
		return []PCComponent{}
	}

	if endIdx > len(components) {
		endIdx = len(components)
	}

	return components[startIdx:endIdx]
}

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(gin.DefaultWriter)

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
		if err := checkError(c, filteredComponents); err != nil {
			return
		}

		sortedComponents := sortComponents(filteredComponents, sortBy)
		if err := checkError(c, sortedComponents); err != nil {
			return
		}

		paginatedComponents := paginateComponents(sortedComponents, pageStr, itemsPerPage)
		if err := checkError(c, paginatedComponents); err != nil {
			return
		}
		log.WithFields(logrus.Fields{
			"TypeFilter":          componentType,
			"BrandFilter":         brandFilter,
			"SortBy":              sortBy,
			"Page":                pageStr,
			"ItemsPerPage":        itemsPerPage,
			"FilteredComponents":  len(filteredComponents),
			"PaginatedComponents": len(paginatedComponents),
		}).Info("Filtered and Paginated Components")
		c.HTML(http.StatusOK, "components.html", gin.H{
			"Components": paginatedComponents,
		})
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{Addr: ":8080", Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	log.Println("Server exiting")

	<-quit
	log.Println("Shutdown signal received, shutting down gracefully...")
	if err := srv.Shutdown(nil); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server gracefully stopped")
}

func getMockPCComponents() []PCComponent {
	return []PCComponent{
		{"CPU", "Intel", "Intel Core i9-10900K", "https://static.shop.kz/upload/resize_cache/iblock/481/450_450_1/154212x1.jpg", 500},
		{"GPU", "NVIDIA", "NVIDIA GeForce RTX 3080", "https://www.nvidia.com/content/dam/en-zz/Solutions/geforce/ampere/rtx-3080-3080ti/geforce-rtx-3080-ti-product-gallery-inline-850-2.jpg", 800},
		{"RAM", "Corsair", "Corsair Vengeance RGB Pro 16GB", "https://static.shop.kz/upload/resize_cache/iblock/75f/jnedijtsz335v3d1uo9uzx6605coss0l/450_450_1/171347o4.jpg", 150},
		{"Motherboard", "ASUS", "ASUS ROG Strix Z490-E Gaming", "https://static.shop.kz/upload/resize_cache/iblock/7fa/9ybjfv07nn75rg1o27ux9q7xdo0cw1yv/450_450_1/171994x1.jpg", 300},
		{"Storage", "Samsung", "Samsung 970 EVO Plus 1TB", "https://static.shop.kz/upload/resize_cache/iblock/113/450_450_1/155791_1.jpg", 200},
		{"Power Supply", "EVGA", "EVGA SuperNOVA 850 G5", "https://static.shop.kz/upload/resize_cache/iblock/68a/450_450_1/158510_01.jpg", 150},
		{"CPU", "AMD", "AMD Ryzen 9 5900X", "https://static.shop.kz/upload/resize_cache/iblock/a05/450_450_1/177555n1.jpg", 550},
		{"GPU", "AMD", "AMD Radeon RX 6800 XT", "https://static.shop.kz/upload/resize_cache/iblock/588/450_450_1/183588n1.jpg", 700},
		{"RAM", "G.Skill", "G.Skill Trident Z Neo 32GB", "https://static.shop.kz/upload/resize_cache/iblock/c3e/450_450_1/175271n1.jpg", 250},
		{"Motherboard", "MSI", "MSI MPG X570 Gaming Pro Carbon WiFi", "https://static.shop.kz/upload/resize_cache/iblock/ba3/450_450_1/177844x1.jpg", 280},
		{"Storage", "Western Digital", "WD Black SN750 NVMe SSD 1TB", "https://static.shop.kz/upload/resize_cache/iblock/79a/450_450_1/179287_1.jpg", 180},
		{"Power Supply", "Corsair", "Corsair RM850x 850W", "https://static.shop.kz/upload/resize_cache/iblock/a42/450_450_1/183238x1.jpg", 160},
	}
}
func checkError(c *gin.Context, data interface{}) error {
	if len(data.([]PCComponent)) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data available"})
		return fmt.Errorf("no data available")
	}
	return nil
}
