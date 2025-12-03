package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
)

type CustomerAuthHandler struct {
	customerAuthService service.CustomerAuthService
}

func NewCustomerAuthHandler(customerAuthService service.CustomerAuthService) *CustomerAuthHandler {
	return &CustomerAuthHandler{
		customerAuthService: customerAuthService,
	}
}

func (h *CustomerAuthHandler) Register(c *fiber.Ctx) error {
	var req dto.CustomerRegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.customerAuthService.Register(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

func (h *CustomerAuthHandler) Login(c *fiber.Ctx) error {
	var req dto.CustomerLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.customerAuthService.Login(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}
