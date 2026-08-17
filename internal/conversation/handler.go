package conversation

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

type CreateConversationRequest struct {
	UserID uint64  `json:"userId"`
	Title  *string `json:"title"`
}

type AddMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateConversationRequest

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

	conversation, err := h.service.CreateConversation(
		req.UserID,
		req.Title,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "failed to create conversation",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    conversation,
	})
}

func (h *Handler) AddMessage(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid conversation id",
		})
	}

	var req AddMessageRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	message, err := h.service.AddMessage(
		id,
		req.Role,
		req.Content,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    message,
	})
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid conversation id",
		})
	}

	conversation, err := h.service.GetConversation(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "conversation not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    conversation,
	})
}
