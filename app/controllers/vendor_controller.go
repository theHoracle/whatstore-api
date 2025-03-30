package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// CreateVendor godoc
// @Summary Create a new vendor account
// @Description Create a new vendor account with an initial store
// @Tags vendors
// @Accept json
// @Produce json
//
//	@Param input body struct {
//		StoreName        string `json:"store_name" binding:"required"`
//		StoreDescription string `json:"store_description" binding:"required"`
//		StoreLogo        string `json:"store_logo" binding:"required"`
//	} true "Vendor details"
//
// @Success 201 {object} models.Vendor
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerAuth
// @Router /vendors [post]
func CreateVendor(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	var input struct {
		StoreName        string `json:"store_name" binding:"required"`
		StoreDescription string `json:"store_description" binding:"required"`
		StoreLogo        string `json:"store_logo" binding:"required"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := c.Locals("user").(*models.User)
	userID := user.ID

	// Check if user already has a vendor account
	var existingVendor models.Vendor
	if err := db.Where("user_id = ?", userID).First(&existingVendor).Error; err == nil {
		// If user already has a vendor account, create a new store for them
		store := models.Store{
			VendorID:    existingVendor.ID,
			Name:        input.StoreName,
			Description: input.StoreDescription,
			StoreLogo:   input.StoreLogo,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := db.Create(&store).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(store)
	}

	// Create new vendor account with initial store
	vendor := models.Vendor{
		UserID:    userID,
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	store := models.Store{
		VendorID:    vendor.ID,
		Name:        input.StoreName,
		Description: input.StoreDescription,
		StoreLogo:   input.StoreLogo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&store).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	vendor.Stores = []models.Store{store}
	return c.Status(fiber.StatusCreated).JSON(vendor)
}

// UpdateVendor godoc
// @Summary Update vendor details
// @Description Update an existing vendor's information
// @Tags vendors
// @Accept json
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} models.Vendor
// @Failure 400 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerAuth
// @Router /vendors/{id} [put]
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

// DeleteVendor godoc
// @Summary Delete a vendor
// @Description Delete a vendor and all associated stores
// @Tags vendors
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 204 "No Content"
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerAuth
// @Router /vendors/{id} [delete]
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

// GetVendor godoc
// @Summary Get vendor details
// @Description Get detailed information about a vendor
// @Tags vendors
// @Produce json
// @Param id path string true "Vendor ID"
// @Success 200 {object} models.Vendor
// @Failure 404 {object} object{error=string}
// @Router /vendors/{id} [get]
func GetVendor(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	id := c.Params("id")
	vendor := new(models.Vendor)

	if err := db.First(&vendor, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vendor not found"})
	}

	return c.JSON(vendor)
}

// GetAllVendors godoc
// @Summary List all vendors
// @Description Get a list of all vendors
// @Tags vendors
// @Produce json
// @Success 200 {array} models.Vendor
// @Failure 500 {object} object{error=string}
// @Router /vendors [get]
func GetAllVendors(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	var vendors []models.Vendor

	if err := db.Find(&vendors).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch vendors"})
	}

	return c.JSON(vendors)
}
