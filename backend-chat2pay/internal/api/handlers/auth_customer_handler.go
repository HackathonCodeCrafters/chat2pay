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

// Register godoc
// @Summary Register Customer
// @Description Registrasi customer baru
// @Tags Auth-Customer
// @Accept json
// @Produce json
// @Param request body dto.CustomerRegisterRequest true "Customer registration data"
// @Success 201 {object} presenter.SuccessResponseSwagger{data=dto.CustomerAuthResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Router /auth/customer/register [post]
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

// Login godoc
// @Summary Login Customer
// @Description Login untuk customer
// @Tags Auth-Customer
// @Accept json
// @Produce json
// @Param request body dto.CustomerLoginRequest true "Login credentials"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.CustomerAuthResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Router /auth/customer/login [post]
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
