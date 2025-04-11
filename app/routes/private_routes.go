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

	// vendors management
	VendorRoutes(api)

	// store management
	StoreRoutes(api)

	// order management
	OrderRoutes(api)

	// Vendor store management
	stores := api.Group("/stores")
	{
		// Store products management
		stores.Post("/:storeId/products", controllers.CreateProduct)
		stores.Put("/:storeId/products/:id", controllers.UpdateProduct)
		stores.Delete("/:storeId/products/:id", controllers.DeleteProduct)

		// Store orders
		stores.Get("/:storeId/orders", controllers.GetStoreOrders)
	}

	// User profile and orders
	users := api.Group("/users")
	{
		users.Get("/me", controllers.GetUserProfile)
		users.Put("/me", controllers.UpdateUserProfile)
	}

	orders := api.Group("/orders")
	{
		orders.Get("/", controllers.GetUserOrders)
		orders.Post("/", controllers.CreateOrder)
		orders.Get("/:id", controllers.GetOrder)
	}

	// Admin dashboard routes
	admin := api.Group("/admin", middleware.AuthMiddleware(db))
	{
		admin.Get("/stats", controllers.GetStats)
		admin.Get("/orders", controllers.GetAllOrders)
		admin.Put("/orders/:id/status", controllers.UpdateOrderStatus)
	}
}
