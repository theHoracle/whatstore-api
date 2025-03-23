package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/middleware"
	"gorm.io/gorm"
)

func PrivateRoutes(a *fiber.App, db *gorm.DB) {
	privateRoutes := a.Group("/api", middleware.AuthMiddleware(db))

	VendorRoutes(privateRoutes)
	StoreRoutes(privateRoutes)
}
