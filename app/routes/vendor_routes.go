package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
)

func VendorRoutes(app fiber.Router) {
	app.Post("/vendors", controllers.CreateVendor)
	app.Put("/vendors/:id", controllers.UpdateVendor)
	app.Delete("/vendors/:id", controllers.DeleteVendor)
	app.Get("/vendors/:id", controllers.GetVendor)
	// app.Get("/vendors", controllers.GetAllVendors)
}
