package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "AVIGO API is running",
		})
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	log.Println("AVIGO API running on :8080")

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}