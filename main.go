package main

import (
	"log"
	"os"

	//	"strings"
	//	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/theHoracle/whatstore-api/db/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("NO ENV FILE FOUND!")
	}

	database.ConnectDB()
	app := fiber.New()

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("Hello whatsapi")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("APP LISTENING ON PORT " + port)
	app.Listen(":" + port)
}
