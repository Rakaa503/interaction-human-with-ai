package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/Rakaa503/AviGo/internal/config"
	appcontext "github.com/Rakaa503/AviGo/internal/context"
	"github.com/Rakaa503/AviGo/internal/conversation"
	"github.com/Rakaa503/AviGo/internal/database"
	"github.com/Rakaa503/AviGo/internal/decision"
	"github.com/Rakaa503/AviGo/internal/interaction"
	"github.com/Rakaa503/AviGo/internal/orchestrator"
	"github.com/Rakaa503/AviGo/internal/response"
)

func main() {
	cfg := config.Load()

	// =========================
	// Database
	// =========================

	db := database.Connect(cfg.DatabaseURL)

	if err := database.Ping(db); err != nil {
		log.Fatalf(
			"database ping failed: %v",
			err,
		)
	}

	// =========================
	// Fiber
	// =========================

	app := fiber.New()

	// =========================
	// Conversation Module
	// =========================

	conversationRepository := conversation.NewRepository(db)

	conversationService := conversation.NewService(
		conversationRepository,
	)

	conversationHandler := conversation.NewHandler(
		conversationService,
	)

	// =========================
	// Interaction Module
	// =========================

	interactionRepository := interaction.NewRepository(db)

	interactionAnalyzer := interaction.NewMLAnalyzer(
		"http://127.0.0.1:8000",
	)

	interactionService := interaction.NewService(
		interactionRepository,
		interactionAnalyzer,
	)

	interactionHandler := interaction.NewHandler(
		interactionService,
	)

	// =========================
	// Context Module
	// =========================

	contextService := appcontext.NewService(
		conversationRepository,
		interactionRepository,
	)

	// =========================
	// Decision Module
	// =========================

	decisionService := decision.NewService()

	// =========================
	// Response Module
	// =========================

	responseService := response.NewService()

	// =========================
	// Orchestrator Module
	// =========================

	orchestratorService := orchestrator.NewService(
		interactionService,
		contextService,
		decisionService,
		responseService,
		conversationService,
	)

	orchestratorHandler := orchestrator.NewHandler(
		orchestratorService,
	)

	// =========================
	// API v1
	// =========================

	api := app.Group("/api/v1")

	// -------------------------
	// Conversation Routes
	// -------------------------

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

	// -------------------------
	// Interaction Routes
	// -------------------------

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

	// -------------------------
	// Chat Route
	// -------------------------

	api.Post(
		"/chat",
		orchestratorHandler.Chat,
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

	// =========================
	// Server
	// =========================

	log.Printf(
		"%s API running on :%s",
		cfg.AppName,
		cfg.AppPort,
	)

	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
