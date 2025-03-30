package controllers

import (
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// GetAllProducts godoc
// @Summary Get all products with pagination
// @Description Get a list of all available products with pagination support
// @Tags products
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} PaginationResponse
// @Router /products [get]
func GetAllProducts(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	page, perPage := paginate(c)

	var products []models.Product
	var total int64

	// Get total count
	db.Model(&models.Product{}).Count(&total)

	// Get paginated products
	err := db.Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&products).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch products",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(PaginationResponse{
		Data:        products,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	})
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products using full-text search
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} PaginationResponse
// @Router /products/search [get]
func SearchProducts(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	query := c.Query("q")
	page, perPage := paginate(c)

	var products []models.Product
	var total int64

	// Using the search_vector column for better performance
	searchQuery := db.Where(
		"search_vector @@ plainto_tsquery('english', ?)",
		query,
	).Order(fmt.Sprintf("ts_rank(search_vector, plainto_tsquery('english', '%s')) DESC", query))

	// Get total count for search results
	searchQuery.Model(&models.Product{}).Count(&total)

	// Get paginated search results
	err := searchQuery.
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&products).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not perform search",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(PaginationResponse{
		Data:        products,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	})
}

// GetAllServices godoc
// @Summary Get all services with pagination
// @Description Get a list of all available services with pagination support
// @Tags services
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} PaginationResponse
// @Router /services [get]
func GetAllServices(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	page, perPage := paginate(c)

	var services []models.Service
	var total int64

	// Get total count
	db.Model(&models.Service{}).Count(&total)

	// Get paginated services
	err := db.Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&services).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch services",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(PaginationResponse{
		Data:        services,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	})
}

// SearchServices godoc
// @Summary Search services
// @Description Search services using full-text search
// @Tags services
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} PaginationResponse
// @Router /services/search [get]
func SearchServices(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	query := c.Query("q")
	page, perPage := paginate(c)

	var services []models.Service
	var total int64

	// Using the search_vector column for better performance
	searchQuery := db.Where(
		"search_vector @@ plainto_tsquery('english', ?)",
		query,
	).Order(fmt.Sprintf("ts_rank(search_vector, plainto_tsquery('english', '%s')) DESC", query))

	// Get total count for search results
	searchQuery.Model(&models.Service{}).Count(&total)

	// Get paginated search results
	err := searchQuery.
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&services).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not perform search",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return c.JSON(PaginationResponse{
		Data:        services,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	})
}
