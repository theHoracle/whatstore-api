package controllers

import (
	"fmt"
	"regexp"
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
// @Param input body models.CreateVendorRequest true "Vendor store details"
// @Success 201 {object} models.Vendor
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /vendors [post]
func CreateVendor(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	// var input models.CreateVendorRequest
	// if err := c.BodyParser(&input); err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// }

	user := c.Locals("user").(*models.User)
	userID := user.ID

	// Check if user already has a vendor account
	var existingVendor models.Vendor
	if err := db.Where("user_id = ?", userID).First(&existingVendor).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already has a vendor account"})
	}

	vendor := models.Vendor{
		UserID:    userID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&vendor).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create vendor"})
	}

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

// CreateStore godoc
// @Summary Create a new store
// @Description Create a new store for the authenticated vendor
// @Tags stores
// @Accept json
// @Produce json
// @Param input body models.CreateStoreRequest true "Store details"
// @Success 201 {object} models.Store
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /store/create [post]
func CreateStore(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)

	var input models.CreateStoreRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Get vendor information for the authenticated user
	var vendor models.Vendor
	if err := db.Where("user_id = ?", user.ID).First(&vendor).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vendor account not found"})
	}

	// Validate the phone number format
	if err := VaidatePhoneNumber(input.StoreWhatsappContact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid phone number format"})
	}
	// Check if the store URL is already taken
	var existingStore models.Store
	if err := db.Where("store_url = ?", input.StoreUrl).First(&existingStore).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Store URL already taken"})
	}

	store := models.Store{
		VendorID:             vendor.ID,
		Name:                 input.StoreName,
		Description:          input.StoreDescription,
		StoreLogo:            input.StoreLogo,
		StoreUrl:             input.StoreUrl,
		StoreAddress:         input.StoreAddress,
		StoreWhatsappContact: input.StoreWhatsappContact,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := db.Create(&store).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(store)
}

// CheckStoreUrlAvailability godoc
// @Summary Check if a store URL is available
// @Description Check if a store URL is already taken
// @Tags stores
// @Accept json
// @Produce json
// @Param url query string true "Store URL to check"
// @Success 200 {object} object{available=boolean}
// @Failure 400 {object} models.ErrorResponse
// @Router /stores/check-url [get]
func CheckStoreUrlAvailability(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	url := c.Query("url")

	if url == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL parameter is required",
		})
	}

	var count int64
	db.Model(&models.Store{}).Where("store_url = ?", url).Count(&count)

	return c.JSON(fiber.Map{
		"available": count == 0,
	})
}

var e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

func VaidatePhoneNumber(phone string) error {
	// Check if the phone number matches the E.164 format
	if !e164Regex.MatchString(phone) {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}
