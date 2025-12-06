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

// Register godoc
// @Summary Register Merchant
// @Description Registrasi merchant baru beserta owner user
// @Tags Auth-Merchant
// @Accept json
// @Produce json
// @Param request body dto.MerchantRegisterRequest true "Merchant registration data"
// @Success 201 {object} presenter.SuccessResponseSwagger{data=dto.MerchantAuthResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Router /auth/merchant/register [post]
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

// Login godoc
// @Summary Login Merchant
// @Description Login untuk merchant user
// @Tags Auth-Merchant
// @Accept json
// @Produce json
// @Param request body dto.MerchantLoginRequest true "Login credentials"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.MerchantAuthResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Router /auth/merchant/login [post]
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
