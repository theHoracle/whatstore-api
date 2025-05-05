package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
	"github.com/theHoracle/whatstore-api/app/middleware"
	"gorm.io/gorm"
)

func PrivateRoutes(app *fiber.App, db *gorm.DB) {
	// API group with version and auth middleware
	api := app.Group("/api/v1", middleware.AuthMiddleware(db))

	// Setup route groups
	VendorRoutes(api)
	StoreRoutes(api)
	OrderRoutes(api)

	// User Management Routes
	users := api.Group("/users")
	{
		users.Get("/me", controllers.GetUserProfile)
		users.Put("/me", controllers.UpdateUserProfile)
	}

	// Admin Routes
	admin := api.Group("/admin", middleware.AuthMiddleware(db))
	{
		admin.Get("/stats", controllers.GetStats)
		admin.Get("/orders", controllers.GetAllOrders)
		admin.Put("/orders/:id/status", controllers.UpdateOrderStatus)
	}
}
