package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/entities"
	"chat2pay/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatRepo repositories.ChatMessageRepository
}

func NewChatHandler(chatRepo repositories.ChatMessageRepository) *ChatHandler {
	return &ChatHandler{chatRepo: chatRepo}
}

// GetHistory godoc
// @Summary Get chat history
// @Description Get chat history for current customer
// @Tags Chat
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /chat/history [get]
func (h *ChatHandler) GetHistory(c *fiber.Ctx) error {
	customerIDVal := c.Locals("customer_id")
	if customerIDVal == nil {
		return c.Status(401).JSON(presenter.ErrorResponse(fiber.NewError(401, "Unauthorized: customer_id not found")))
	}
	customerID := customerIDVal.(string)

	messages, err := h.chatRepo.FindByCustomerID(c.Context(), customerID, 100)
	if err != nil {
		return c.Status(500).JSON(presenter.ErrorResponse(err))
	}

	return c.JSON(presenter.SuccessResponse(dto.ToChatHistoryResponse(messages)))
}

// SaveMessage godoc
// @Summary Save chat message
// @Description Save a chat message
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.ChatMessageRequest true "Message data"
// @Success 201 {object} presenter.SuccessResponseSwagger
// @Router /chat/messages [post]
func (h *ChatHandler) SaveMessage(c *fiber.Ctx) error {
	customerIDVal := c.Locals("customer_id")
	if customerIDVal == nil {
		return c.Status(401).JSON(presenter.ErrorResponse(fiber.NewError(401, "Unauthorized: customer_id not found")))
	}
	customerID := customerIDVal.(string)

	var req dto.ChatMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	msg := &entities.ChatMessage{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		Role:       req.Role,
		Content:    req.Content,
		Products:   req.Products,
	}

	if err := h.chatRepo.Create(c.Context(), msg); err != nil {
		return c.Status(500).JSON(presenter.ErrorResponse(err))
	}

	return c.Status(201).JSON(presenter.SuccessResponse(dto.ToChatMessageResponse(msg)))
}

// ClearHistory godoc
// @Summary Clear chat history
// @Description Clear all chat history for current customer
// @Tags Chat
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /chat/history [delete]
func (h *ChatHandler) ClearHistory(c *fiber.Ctx) error {
	customerIDVal := c.Locals("customer_id")
	if customerIDVal == nil {
		return c.Status(401).JSON(presenter.ErrorResponse(fiber.NewError(401, "Unauthorized: customer_id not found")))
	}
	customerID := customerIDVal.(string)

	if err := h.chatRepo.DeleteByCustomerID(c.Context(), customerID); err != nil {
		return c.Status(500).JSON(presenter.ErrorResponse(err))
	}

	return c.JSON(presenter.SuccessResponse(map[string]string{"message": "Chat history cleared"}))
}
