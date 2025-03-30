package main

import (
	"log"
	"os"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
	swagger "github.com/swaggo/fiber-swagger"
	"github.com/theHoracle/whatstore-api/app/handlers"
	"github.com/theHoracle/whatstore-api/app/routes"
	"github.com/theHoracle/whatstore-api/db/database"
	_ "github.com/theHoracle/whatstore-api/docs" // This will import the generated docs
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("NO ENV FILE FOUND!")
	}

	database.ConnectDB()

	// init clerk
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))
	if os.Getenv("CLERK_SECRET_KEY") == "" {
		log.Fatal("Clerk secret key not set")
	}

	clerkSigningSecret := os.Getenv("CLERK_SIGNING_SECRET")
	log.Println("Clerk signing secret: " + clerkSigningSecret)
	if clerkSigningSecret == "" {
		log.Fatal("Clerk signing secret not set")
	}

	app := fiber.New()

	// Add rate limiter middleware
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			// skip if non localhost
			return c.IP() == "127.0.0.1"
		},
		Max:        100,             // max number of requests
		Expiration: 1 * time.Minute, // per 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // use IP as key for rate limiting
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	// Global middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", database.DB.Db)
		return c.Next()
	})

	// Documentation routes
	app.Get("/swagger/*", swagger.WrapHandler)

	// Webhooks
	webhooks := app.Group("/webhooks")
	{
		webhooks.Post("/clerk", handlers.ClerkWebhookHandler(database.DB.Db, clerkSigningSecret))
	}

	// Setup API routes
	routes.PublicRoutes(app)
	routes.PrivateRoutes(app, database.DB.Db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("APP LISTENING ON PORT " + port)
	app.Listen(":" + port)
}
