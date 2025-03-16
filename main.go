package main

import (
	"log"
	"os"

	//	"strings"
	//	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/theHoracle/whatstore-api/app/routes"
	"github.com/theHoracle/whatstore-api/db/database"
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

	app := fiber.New()

	routes.PublicRoutes(app)
	routes.VendorRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("APP LISTENING ON PORT " + port)
	app.Listen(":" + port)
}
