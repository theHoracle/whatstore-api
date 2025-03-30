package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// GetStats godoc
// @Summary Get admin dashboard stats
// @Description Get statistics for admin dashboard
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.DashboardStats
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/stats [get]
func GetStats(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	var stats models.DashboardStats

	// Get total orders
	db.Model(&models.Order{}).Count(&stats.TotalOrders)

	// Get total products
	db.Model(&models.Product{}).Count(&stats.TotalProducts)

	// Get total users
	db.Model(&models.User{}).Count(&stats.TotalUsers)

	return c.JSON(stats)
}

// GetAllOrders godoc
// @Summary Get all orders (admin)
// @Description Get all orders for admin dashboard
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} PaginationResponse{data=[]models.Order}
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/orders [get]
func GetAllOrders(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	page, perPage := paginate(c)

	var orders []models.Order
	var total int64

	db.Model(&models.Order{}).Count(&total)
	if err := db.Offset((page - 1) * perPage).Limit(perPage).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(NewPaginationResponse(orders, total, page, perPage))
}

// UpdateOrderStatus godoc
// @Summary Update order status (admin)
// @Description Update the status of an order
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param status body models.UpdateOrderStatusRequest true "New order status"
// @Success 200 {object} models.Order
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/admin/orders/{id}/status [put]
func UpdateOrderStatusAdmin(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	orderID := c.Params("id")

	var statusRequest models.UpdateOrderStatusRequest
	if err := c.BodyParser(&statusRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var order models.Order
	if err := db.First(&order, orderID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	order.Status = models.OrderStatus(statusRequest.Status)
	if err := db.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	return c.JSON(order)
}
