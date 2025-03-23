package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
)

func PublicRoutes(a *fiber.App) {
	publicRoutes := a.Group("/api")

	publicRoutes.Get("/vendors", controllers.GetAllVendors)
	publicRoutes.Get("/vendors/:id", controllers.GetVendor)

}
