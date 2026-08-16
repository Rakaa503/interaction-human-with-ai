package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/Rakaa503/AviGo/internal/config"
	"github.com/Rakaa503/AviGo/internal/database"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg.DatabaseURL)

	if err := database.Ping(db); err != nil {
		log.Fatalf(
			"database ping failed: %v",
			err,
		)
	}

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

	log.Printf(
		"%s API running on :%s",
		cfg.AppName,
		cfg.AppPort,
	)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
