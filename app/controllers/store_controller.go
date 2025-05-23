package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// UpdateStore godoc
// @Summary Update store details
// @Description Update an existing store's information
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param input body object true "Store update information" SchemaExample({"name":"Store Name","description":"Store Description","store_logo":"http://example.com/image.jpg"})
// @Success 200 {object} models.Store
// @Failure 400 {object} object{error=string}
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerAuth
// @Router /stores/{id} [put]
func UpdateStore(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	storeID := c.Params("id")
	vendor := c.Locals("user").(*models.User)

	var store models.Store
	if err := db.First(&store, storeID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Store not found"})
	}

	// Verify store belongs to vendor
	var vendorStore models.Vendor
	if err := db.First(&vendorStore, store.VendorID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vendor not found"})
	}

	if vendorStore.UserID != vendor.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized to modify this store"})
	}

	var input struct {
		Name                 string `json:"name"`
		Description          string `json:"description"`
		StoreLogo            string `json:"store_logo"`
		StoreUrl             string `json:"store_url"`
		StoreAddress         string `json:"store_address"`
		StoreWhatsappContact string `json:"store_whatsapp_contact"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Update all fields
	store.Name = input.Name
	store.Description = input.Description
	store.StoreLogo = input.StoreLogo
	store.StoreUrl = input.StoreUrl
	store.StoreAddress = input.StoreAddress
	store.StoreWhatsappContact = input.StoreWhatsappContact
	store.UpdatedAt = time.Now()

	if err := db.Save(&store).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update store"})
	}

	return c.JSON(store)
}

// DeleteStore godoc
// @Summary Delete a store
// @Description Delete a store and all associated data
// @Tags stores
// @Produce json
// @Param id path string true "Store ID"
// @Success 204 "No Content"
// @Failure 403 {object} object{error=string}
// @Failure 404 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Security BearerAuth
// @Router /stores/{id} [delete]
func DeleteStore(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	storeID := c.Params("id")
	vendor := c.Locals("user").(*models.User)

	var store models.Store
	if err := db.First(&store, storeID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Store not found"})
	}

	// Verify ownership
	var vendorStore models.Vendor
	if err := db.First(&vendorStore, store.VendorID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Vendor not found"})
	}

	if vendorStore.UserID != vendor.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized to delete this store"})
	}

	// Start transaction
	tx := db.Begin()

	// Delete all associated products and services
	if err := tx.Where("store_id = ?", store.ID).Delete(&models.Product{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete products"})
	}

	if err := tx.Delete(&store).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not delete store"})
	}

	// if vendor has no store set vendor.isActive to false
	var vendorStores []models.Store
	if err := tx.Where("vendor_id = ?", vendorStore.ID).Find(&vendorStores).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch vendor stores"})
	}

	if len(vendorStores) == 0 {
		vendorStore.IsActive = false
		if err := tx.Save(&vendorStore).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update vendor"})
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not commit transaction"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetStore godoc
// @Summary Get store details
// @Description Get detailed information about a store including its products and services
// @Tags stores
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} models.Store
// @Failure 404 {object} object{error=string}
// @Router /stores/{id} [get]
func GetStore(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	storeID := c.Params("id")

	var store models.Store
	if err := db.Preload("Products").Preload("Services").First(&store, storeID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Store not found"})
	}

	return c.JSON(store)
}

// GetAllStores godoc
// @Summary List all stores for a vendor
// @Description Get a list of all stores belonging to a specific vendor
// @Tags stores
// @Produce json
// @Param vendorId path string true "Vendor ID"
// @Success 200 {array} models.Store
// @Failure 500 {object} object{error=string}
// @Router /vendors/{vendorId}/stores [get]
func GetAllStores(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	vendorID := c.Params("vendorId")

	var stores []models.Store
	if err := db.Where("vendor_id = ?", vendorID).Find(&stores).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not fetch stores"})
	}

	return c.JSON(stores)
}
