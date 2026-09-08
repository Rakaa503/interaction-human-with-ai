package knowledge

import "github.com/gofiber/fiber/v3"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type CreateDocumentRequest struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	URL      *string `json:"url"`
	Source   *string `json:"source"`
	Category *string `json:"category"`
}

func (h *Handler) CreateDocument(c fiber.Ctx) error {
	var request CreateDocumentRequest

	if err := c.Bind().Body(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	document, err := h.service.AddDocument(
		request.Title,
		request.Content,
		request.URL,
		request.Source,
		request.Category,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"data":    document,
	})
}

func (h *Handler) GetDocuments(c fiber.Ctx) error {
	documents, err := h.service.GetAllDocuments()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    documents,
	})
}
