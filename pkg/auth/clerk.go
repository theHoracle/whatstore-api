package auth

import (
	"errors"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
)

// GetUserID retrieves the user ID from the Fiber context
func GetUserID(c *fiber.Ctx) (string, error) {
	user, ok := c.Locals("user").(map[string]interface{})
	if !ok {
		return "", errors.New("user claims not found in context")
	}

	sub, ok := user["sub"].(string)
	if !ok {
		return "", errors.New("sub claim not found or not a string")
	}

	return sub, nil
}

// ClerkAuthMiddleware verifies the JWT from the Authorization header
func ClerkAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Split the header to extract the token (expected format: "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header",
			})
		}

		token := parts[1]

		// Verify the token using Clerk's JWT verification
		verifiedToken, err := jwt.Verify(c.Context(), &jwt.VerifyParams{
			Token: token,
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Extract claims and store them in the context
		c.Locals("user", verifiedToken.Claims)
		return c.Next()
	}
}
