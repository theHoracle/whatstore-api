package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/controllers"
)

func OrderRoutes(app fiber.Router) {
	orders := app.Group("/orders")

	orders.Post("/", controllers.CreateOrder)
	orders.Get("/", controllers.GetUserOrders)
	orders.Get("/store/:storeId", controllers.GetStoreOrders)
	orders.Put("/:id/status", controllers.UpdateOrderStatus)
}
