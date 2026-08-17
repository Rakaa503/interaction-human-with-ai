package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/Rakaa503/AviGo/internal/config"
	"github.com/Rakaa503/AviGo/internal/conversation"
	"github.com/Rakaa503/AviGo/internal/database"
	"github.com/Rakaa503/AviGo/internal/interaction"
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

	// =========================
	// Conversation module
	// =========================

	conversationRepository := conversation.NewRepository(db)

	conversationService := conversation.NewService(
		conversationRepository,
	)

	conversationHandler := conversation.NewHandler(
		conversationService,
	)

	// =========================
	// Interaction module
	// =========================

	interactionRepository := interaction.NewRepository(db)

	interactionAnalyzer := interaction.NewRuleBasedAnalyzer()

	interactionService := interaction.NewService(
		interactionRepository,
		interactionAnalyzer,
	)

	interactionHandler := interaction.NewHandler(
		interactionService,
	)

	// =========================
	// API
	// =========================

	api := app.Group("/api/v1")

	// Conversation routes

	api.Post(
		"/conversations",
		conversationHandler.Create,
	)

	api.Get(
		"/conversations/:id",
		conversationHandler.Get,
	)

	api.Post(
		"/conversations/:id/messages",
		conversationHandler.AddMessage,
	)

	// Interaction routes

	api.Post(
		"/interactions",
		interactionHandler.Create,
	)

	api.Get(
		"/interactions/:id",
		interactionHandler.GetByID,
	)

	api.Get(
		"/conversations/:id/interactions",
		interactionHandler.GetByConversationID,
	)

	// =========================
	// Root
	// =========================

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "AVIGO API is running",
		})
	})

	// =========================
	// Health
	// =========================

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
