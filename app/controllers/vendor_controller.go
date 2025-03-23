package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// CreateVendor handles the creation of a new vendor
func CreateVendor(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	var input struct {
		StoreName        string `json:"store_name" binding:"required"`
		StoreDescription string `json:"store_description" binding:"required"`
		ImageURL         string `json:"image_url" binding:"required"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := c.Locals("user").(*models.User)
	userID := user.ID

	// Check if user already has a vendor account
	var existingVendor models.Vendor
	if err := db.Where("user_id = ?", userID).First(&existingVendor).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already has a vendor account"})
	}

	vendor := models.Vendor{
		UserID:           userID,
		StoreName:        input.StoreName,
		StoreDescription: input.StoreDescription,
		ImageURL:         input.ImageURL,
		IsActive:         false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(vendor)
}

// UpdateVendor handles updating an existing vendor
func UpdateVendor(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	id := c.Params("id")
	// currentUserID := c.Locals("user").(*models.User).ID

	vendor := new(models.Vendor)

	if err := db.First(&vendor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vendor not found"})
	}

	if err := c.BodyParser(vendor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := db.Save(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update vendor"})
	}

	return c.JSON(vendor)
}

// DeleteVendor handles deleting a vendor
func DeleteVendor(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	id := c.Params("id")
	vendor := new(models.Vendor)

	if err := db.First(&vendor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vendor not found"})
	}

	if err := db.Delete(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot delete vendor"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetVendor handles fetching a vendor by ID
func GetVendor(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	id := c.Params("id")
	vendor := new(models.Vendor)

	if err := db.First(&vendor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vendor not found"})
	}

	return c.JSON(vendor)
}

// GetAllVendors handles fetching all vendors with pagination
func GetAllVendors(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	var vendors []models.Vendor

	if err := db.Find(&vendors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch vendors"})
	}

	return c.JSON(vendors)
}
