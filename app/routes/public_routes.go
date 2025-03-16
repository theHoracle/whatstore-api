package routes

import "github.com/gofiber/fiber/v2"

func PublicRoutes(a *fiber.App) {
	publicRoutes := a.Group("/api")

	publicRoutes.Get("/vendors", func(c *fiber.Ctx) error {
		return c.SendString("All Vendors with pagination")
	})
}
