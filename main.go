package main

import (
	"github-praiser/handlers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	app := fiber.New()

	app.Post("/praising", handlers.HandlePraising)

	app.Listen(":3000")
}
