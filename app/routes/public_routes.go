package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
)

func PublicRoutes(a *fiber.App) {
	publicRoutes := a.Group("/api")

	publicRoutes.Get("/vendors", controllers.GetAllVendors)
	publicRoutes.Get("/vendors/:id", controllers.GetVendor)

	// Product routes
	publicRoutes.Get("/products", controllers.GetAllProducts)
	publicRoutes.Get("/products/search", controllers.SearchProducts)

	// Service routes
	publicRoutes.Get("/services", controllers.GetAllServices)
	publicRoutes.Get("/services/search", controllers.SearchServices)
}
