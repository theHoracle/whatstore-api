package main

import (
	"log"
	"os"
	//	"strings"
	//	"time"

	//	"github.com/go-pkgz/auth"
	//	"github.com/go-pkgz/auth/token"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("NO ENV FILE FOUND!")
	}
	//	options := auth.Opts{
	//	SecretReader: token.SecretFunc(func(id string) (string, error) {
	//			return "secret", nil // get secret from env
	//		}),
	//		TokenDuration:  time.Minute * 5,
	//		CookieDuration: time.Hour * 24,
	//		Issuer:         "MY APP",
	//		URL:            "",
	//		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
	// allow only dev_* names
	//			return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
	//		}),
	//	}

	//	authService := auth.NewService(options)
	//	authService.AddProvider("google", "key", "secret")

	//	m := authService.Middleware()

	app := fiber.New()

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("Hello whatsapi")
	})

	//	app.Use(m.Auth).Get("/auth", func(c *fiber.Ctx) error {
	//		return c.SendString("Hello, welcome here")
	//	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("APP LISTENING ON PORT 3000")
	app.Listen(":" + port)
}
