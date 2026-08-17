package interaction

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type createInteractionRequest struct {
	UserID         uint64 `json:"userId"`
	ConversationID uint64 `json:"conversationId"`
	Input          string `json:"input"`
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req createInteractionRequest

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

	interaction, err := h.service.Process(
		req.UserID,
		req.ConversationID,
		req.Input,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to process interaction",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    interaction,
	})
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		64,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid interaction id",
		})
	}

	interaction, err := h.service.GetByID(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "interaction not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    interaction,
	})
}

func (h *Handler) GetByConversationID(c fiber.Ctx) error {
	id, err := strconv.ParseUint(
		c.Params("id"),
		10,
		64,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid conversation id",
		})
	}

	interactions, err := h.service.GetByConversationID(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to get interactions",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    interactions,
	})
}
