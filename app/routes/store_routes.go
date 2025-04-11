package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
)

func StoreRoutes(app fiber.Router) {
	stores := app.Group("/stores")

	// Store creation
	stores.Post("/create", controllers.CreateStore) // Updated path without :id

	// Store management
	stores.Put("/:id", controllers.UpdateStore)
	stores.Delete("/:id", controllers.DeleteStore)
	stores.Get("/:id", controllers.GetStore)
	stores.Get("/vendor/:vendorId", controllers.GetAllStores)
}
