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
			// Define a struct that matches Clerk's JSON structure
			var clerkUser struct {
				ID             string `json:"id"`
				FirstName      string `json:"first_name"`
				LastName       string `json:"last_name"`
				Username       string `json:"username"`
				ImageURL       string `json:"profile_image_url"`
				EmailAddresses []struct {
					EmailAddress string `json:"email_address"`
				} `json:"email_addresses"`
			}

			// Debug print the raw data
			log.Printf("Raw webhook data: %s", string(event.Data))

			if err := json.Unmarshal(event.Data, &clerkUser); err != nil {
				log.Printf("Error unmarshaling user data: %v", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid user data: " + err.Error(),
				})
			}

			// Debug print the parsed data
			log.Printf("Parsed clerk user: %+v", clerkUser)

			// Extract email from the first email address if available
			email := ""
			if len(clerkUser.EmailAddresses) > 0 {
				email = clerkUser.EmailAddresses[0].EmailAddress
			}

			// Validate required fields
			if email == "" {
				log.Printf("No email address found in webhook data")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "No email address provided",
				})
			}

			if clerkUser.Username == "" {
				log.Printf("No username found in webhook data")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "No username provided",
				})
			}

			// Create the user in the database
			newUser := models.User{
				ClerkID:   clerkUser.ID,
				Name:      clerkUser.FirstName + " " + clerkUser.LastName,
				Email:     email,
				Username:  clerkUser.Username,
				AvatarURL: clerkUser.ImageURL,
			}

			// Debug print the user we're about to create
			log.Printf("Creating user with data: %+v", newUser)

			if err := db.Create(&newUser).Error; err != nil {
				log.Printf("Failed to create user: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create user: " + err.Error(),
				})
			}

		case "user.updated":
			// Update the user.updated case with similar structure
			var clerkUser struct {
				ID             string `json:"id"`
				FirstName      string `json:"first_name"`
				LastName       string `json:"last_name"`
				Username       string `json:"username"`
				ImageURL       string `json:"profile_image_url"`
				EmailAddresses []struct {
					EmailAddress string `json:"email_address"`
				} `json:"email_addresses"`
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
			user.Email = clerkUser.EmailAddresses[0].EmailAddress
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
