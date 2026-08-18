package orchestrator

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type ChatRequest struct {
	UserID         uint64 `json:"userId"`
	ConversationID uint64 `json:"conversationId"`
	Input          string `json:"input"`
}

func (h *Handler) Chat(c fiber.Ctx) error {
	var req ChatRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if req.UserID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "userId is required",
		})
	}

	if req.ConversationID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "conversationId is required",
		})
	}

	if req.Input == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "input is required",
		})
	}

	result, err := h.service.Process(
		req.UserID,
		req.ConversationID,
		req.Input,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to process chat",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}
