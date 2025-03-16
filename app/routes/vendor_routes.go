package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/pkg/auth"
)

func VendorRoutes(app *fiber.App) {
	vendorsGroup := app.Group("/api/vendor", auth.ClerkAuthMiddleware())

	vendorsGroup.Get("/create", func(c *fiber.Ctx) error {
		return c.SendString("Protected vendor app")
	})
}
