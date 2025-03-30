package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"github.com/theHoracle/whatstore-api/db/database"
	"gorm.io/gorm"
)

// GetUserOrders godoc
// @Summary Get user orders
// @Description Get all orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} PaginationResponse{data=[]models.Order}
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/orders [get]
func GetUserOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	db := c.Locals("db").(*gorm.DB)
	page, perPage := paginate(c)

	var orders []models.Order
	var total int64

	db.Model(&models.Order{}).Where("user_id = ?", user.ID).Count(&total)
	query := db.Where("user_id = ?", user.ID).Offset((page - 1) * perPage).Limit(perPage)
	if err := query.Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(NewPaginationResponse(orders, total, page, perPage))
}

// CreateOrder godoc
// @Summary Create new order
// @Description Create a new order for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body models.CreateOrderRequest true "Order creation data"
// @Success 201 {object} models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/orders [post]
func CreateOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	db := c.Locals("db").(*gorm.DB)

	var orderRequest models.CreateOrderRequest
	if err := c.BodyParser(&orderRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Start a transaction
	tx := db.Begin()

	order := models.Order{
		UserID: user.ID,
		Status: models.OrderStatusPending,
	}

	// First create the order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	// Get all products and create order items
	for _, item := range orderRequest.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Product not found: " + strconv.Itoa(int(item.ProductID)),
			})
		}

		// Set store ID from first product if not set
		if order.StoreID == 0 {
			order.StoreID = product.StoreID
		} else if order.StoreID != product.StoreID {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "All products must be from the same store",
			})
		}

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Store current price
		}

		orderItems = append(orderItems, orderItem)
		totalAmount += product.Price * float64(item.Quantity)
	}

	// Update order with total and store ID
	order.TotalAmount = totalAmount
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Create all order items
	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order items",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Fetch complete order with items
	if err := db.Preload("Items.Product").First(&order, order.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch created order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// GetOrder godoc
// @Summary Get order details
// @Description Get details of a specific order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/orders/{id} [get]
func GetOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	db := c.Locals("db").(*gorm.DB)
	orderID := c.Params("id")

	var order models.Order
	if err := db.Where("id = ? AND user_id = ?", orderID, user.ID).First(&order).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.JSON(order)
}

// TODO: Use to validate order products
// func validateOrderProducts(order *models.Order) error {
// 	if len(order.Items) == 0 {
// 		return errors.New("order has no products")
// 	}

// 	storeID := order.Items[0].StoreID
// 	for _, product := range order.Items {
// 		if product.StoreID != storeID {
// 			return errors.New("all products must belong to the same store")
// 		}
// 	}

// 	return nil
// }

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
