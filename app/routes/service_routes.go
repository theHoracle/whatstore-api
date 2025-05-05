package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
)

func ServiceRoutes(app fiber.Router) {
	services := app.Group("/stores/:storeId/services")

	services.Post("/", controllers.CreateService)
	services.Put("/:id", controllers.UpdateService)
	services.Delete("/:id", controllers.DeleteService)
	services.Get("/", controllers.GetStoreServices)
}
