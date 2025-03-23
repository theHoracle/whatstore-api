package handlers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	svix "github.com/svix/svix-webhooks/go"
	"github.com/theHoracle/whatstore-api/app/models"
	"gorm.io/gorm"
)

// ClerkWebhookHandler handles incoming Clerk webhooks
func ClerkWebhookHandler(db *gorm.DB, signingSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract raw body and Svix headers for verification
		body := c.Body()
		headers := c.GetReqHeaders()
		svixID := c.Get("svix-id")
		svixTimestamp := c.Get("svix-timestamp")
		svixSignature := c.Get("svix-signature")

		// Check for missing headers
		if svixID == "" || svixTimestamp == "" || svixSignature == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing Svix headers",
			})
		}
		wh, err := svix.NewWebhook(signingSecret)
		if err != nil {
			log.Printf("Failed to create webhook: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create webhook",
			})
		}

		// Parse the payload
		// var payload map[string]interface{}
		// if err := json.Unmarshal(body, &payload); err != nil {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"error": "Invalid payload",
		// 	})
		// }

		// Verify the webhook signature
		err = wh.Verify(body, headers)
		if err != nil {
			log.Printf("Webhook verification failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid webhook signature",
			})
		}

		// Parse the event payload
		var event struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(body, &event); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid webhook payload",
			})
		}

		// Handle each event type
		switch event.Type {
		case "user.created":
			var clerkUser struct {
				ID       string `json:"id"`
				Email    string `json:"email_addresses[0].email_address"`
				Username string `json:"username"`
			}
			if err := json.Unmarshal(event.Data, &clerkUser); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid user data",
				})
			}
			// Create the user in the database
			newUser := models.User{
				ClerkID:  clerkUser.ID,
				Email:    clerkUser.Email,
				Username: clerkUser.Username,
			}
			if err := db.Create(&newUser).Error; err != nil {
				log.Printf("Failed to create user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create user",
				})
			}

		case "user.updated":
			var clerkUser struct {
				ID    string `json:"id"`
				Email string `json:"email_addresses[0].email_address"`
			}
			if err := json.Unmarshal(event.Data, &clerkUser); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid user data",
				})
			}
			// Update the user in the database
			var user models.User
			if err := db.Where("clerk_id = ?", clerkUser.ID).First(&user).Error; err != nil {
				log.Printf("User not found: %v", err)
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			user.Email = clerkUser.Email
			if err := db.Save(&user).Error; err != nil {
				log.Printf("Failed to update user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to update user",
				})
			}

		case "user.deleted":
			var clerkUser struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(event.Data, &clerkUser); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid user data",
				})
			}
			// Delete the user from the database
			if err := db.Where("clerk_id = ?", clerkUser.ID).Delete(&models.User{}).Error; err != nil {
				log.Printf("Failed to delete user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to delete user",
				})
			}

		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}

		// Respond with success
		return c.SendStatus(fiber.StatusOK)
	}
}
