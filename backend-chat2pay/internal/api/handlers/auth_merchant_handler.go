package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
)

type MerchantAuthHandler struct {
	merchantAuthService service.MerchantAuthService
}

func NewMerchantAuthHandler(merchantAuthService service.MerchantAuthService) *MerchantAuthHandler {
	return &MerchantAuthHandler{
		merchantAuthService: merchantAuthService,
	}
}

func (h *MerchantAuthHandler) Register(c *fiber.Ctx) error {
	var req dto.MerchantRegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.merchantAuthService.Register(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

func (h *MerchantAuthHandler) Login(c *fiber.Ctx) error {
	var req dto.MerchantLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.merchantAuthService.Login(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}
