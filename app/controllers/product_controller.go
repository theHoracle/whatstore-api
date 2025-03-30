package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// Helper function to validate store ownership
func validateStoreOwnership(db *gorm.DB, storeID uint, vendorID uint) error {
	var store models.Store
	if err := db.Where("id = ? AND vendor_id = ?", storeID, vendorID).First(&store).Error; err != nil {
		return errors.New("store not found or not authorized")
	}
	return nil
}

// GetProduct godoc
// @Summary Get a single product
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/products/{id} [get]
func GetProduct(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	id := c.Params("id")

	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product for a specific store
// @Tags store-products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param storeId path string true "Store ID"
// @Param product body models.CreateProductRequest true "Product creation data"
// @Success 201 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/v1/stores/{storeId}/products [post]
func CreateProduct(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	storeID, err := c.ParamsInt("storeId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid store ID",
		})
	}

	// Validate store ownership
	if err := validateStoreOwnership(db, uint(storeID), user.Vendor.ID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var productRequest models.CreateProductRequest
	if err := c.BodyParser(&productRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	product := models.Product{
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		StoreID:     uint(storeID),
		ImageURL:    productRequest.ImageURL,
		Stock:       productRequest.Stock,
		Category:    productRequest.Category,
	}

	if err := db.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product (admin only)
// @Tags admin-products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param product body models.UpdateProductRequest true "Product update data"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/admin/products/{id} [put]
func UpdateProduct(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	productID := c.Params("id")

	var product models.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Validate store ownership
	if err := validateStoreOwnership(db, product.StoreID, user.Vendor.ID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var updateData models.UpdateProductRequest
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := db.Model(&product).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete an existing product (admin only)
// @Tags admin-products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/admin/products/{id} [delete]
func DeleteProduct(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	productID := c.Params("id")

	var product models.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Validate store ownership
	if err := validateStoreOwnership(db, product.StoreID, user.Vendor.ID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
