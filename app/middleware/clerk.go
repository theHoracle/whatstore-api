package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// AuthMiddleware verifies Clerk JWT tokens and attaches the user to the context
func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Ensure it’s in "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header. Expected: Bearer <token>",
			})
		}
		token := parts[1]

		// Verify the token using Clerk SDK
		claims, err := jwt.Verify(context.Background(), &jwt.VerifyParams{
			Token: token,
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Get the user ID (Clerk’s "sub" claim)
		userID := claims.Subject
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token missing user ID",
			})
		}

		// Fetch the user from the database
		var user models.User
		if err := db.Where("clerk_id = ?", userID).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "User not found in database",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		// Attach the user to the context
		c.Locals("user", &user)

		// Proceed to the next handler
		return c.Next()
	}
}
