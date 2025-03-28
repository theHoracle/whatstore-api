package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"github.com/theHoracle/whatstore-api/db/database"
)

func CreateOrder(c *fiber.Ctx) error {
	// Get current user from context
	user := c.Locals("user").(models.User)

	order := new(models.Order)
	if err := c.BodyParser(order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Set user ID
	order.UserID = user.ID

	// Validate all products belong to same store
	if err := validateOrderProducts(order); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Create order in pending state
	db := database.DB.Db
	if err := db.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order"})
	}

	return c.JSON(order)
}

func validateOrderProducts(order *models.Order) error {
	if len(order.Items) == 0 {
		return errors.New("order has no products")
	}

	storeID := order.Items[0].Product.StoreID
	for _, product := range order.Items {
		if product.Product.StoreID != storeID {
			return errors.New("all products must belong to the same store")
		}
	}

	return nil
}

func GetUserOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	var orders []models.Order
	if err := database.DB.Db.Where("user_id = ?", user.ID).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	return c.JSON(orders)
}

func GetStoreOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	storeID := c.Params("storeId")

	// Verify user owns the store
	var store models.Store
	if err := database.DB.Db.Where("id = ? AND vendor_id = ?", storeID, user.Vendor.ID).First(&store).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Not authorized"})
	}

	var orders []models.Order
	if err := database.DB.Db.Where("store_id = ?", storeID).Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	return c.JSON(orders)
}

func UpdateOrderStatus(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)
	orderID := c.Params("id")

	var order models.Order
	if err := database.DB.Db.Where("id = ? AND user_id = ?", orderID, user.ID).First(&order).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{"error": "Not authorized"})
	}

	status := c.Query("status")
	if status != "" {
		order.Status = models.OrderStatus(status)
		// Here we can add payment/escrow release logic when implementing payments
		if err := database.DB.Db.Save(&order).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update order"})
		}
	}

	return c.JSON(order)
}
