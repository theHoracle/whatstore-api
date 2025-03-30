package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
)

func PublicRoutes(app *fiber.App) {
	// API group with version
	api := app.Group("/api/v1")

	// Products endpoints
	products := api.Group("/products")
	{
		products.Get("/", controllers.GetAllProducts)       // List all products
		products.Get("/:id", controllers.GetProduct)        // Get single product
		products.Get("/search", controllers.SearchProducts) // Search products
	}

	// Categories endpoints
	// categories := api.Group("/categories")
	// {
	// 	categories.Get("/", controllers.GetCategories)   // List all categories
	// 	categories.Get("/:id", controllers.GetCategory)  // Get single category
	// }

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
