package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service for a specific store
// @Tags store-services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param storeId path string true "Store ID"
// @Param service body models.CreateServiceRequest true "Service creation data"
// @Success 201 {object} models.Service
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/v1/stores/{storeId}/services [post]
func CreateService(c *fiber.Ctx) error {
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

	var serviceRequest models.CreateServiceRequest
	if err := c.BodyParser(&serviceRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	service := models.Service{
		Name:        serviceRequest.Name,
		Description: serviceRequest.Description,
		Rate:        serviceRequest.Rate,
		Currency:    serviceRequest.Currency,
		StoreID:     uint(storeID),
		ImageURL:    serviceRequest.ImageURL,
	}

	if err := db.Create(&service).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create service",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(service)
}

// UpdateService godoc
// @Summary Update a service
// @Description Update an existing service
// @Tags store-services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param storeId path string true "Store ID"
// @Param id path string true "Service ID"
// @Param service body models.UpdateServiceRequest true "Service update data"
// @Success 200 {object} models.Service
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/stores/{storeId}/services/{id} [put]
func UpdateService(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	serviceID := c.Params("id")

	var service models.Service
	if err := db.Where("id = ?", serviceID).First(&service).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	// Validate store ownership
	if err := validateStoreOwnership(db, service.StoreID, user.Vendor.ID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var updateData models.UpdateServiceRequest
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := db.Model(&service).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update service",
		})
	}

	return c.JSON(service)
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete an existing service
// @Tags store-services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param storeId path string true "Store ID"
// @Param id path string true "Service ID"
// @Success 200 {object} object{message=string}
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/stores/{storeId}/services/{id} [delete]
func DeleteService(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	serviceID := c.Params("id")

	var service models.Service
	if err := db.Where("id = ?", serviceID).First(&service).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Service not found",
		})
	}

	// Validate store ownership
	if err := validateStoreOwnership(db, service.StoreID, user.Vendor.ID); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := db.Delete(&service).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete service",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Service deleted successfully",
	})
}

// GetStoreServices godoc
// @Summary Get all services for a store
// @Description Get all services for a specific store
// @Tags store-services
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param storeId path string true "Store ID"
// @Success 200 {array} models.Service
// @Failure 403 {object} models.ErrorResponse
// @Router /api/v1/stores/{storeId}/services [get]
func GetStoreServices(c *fiber.Ctx) error {
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

	var services []models.Service
	if err := db.Where("store_id = ?", storeID).Find(&services).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch services",
		})
	}

	return c.JSON(services)
}
