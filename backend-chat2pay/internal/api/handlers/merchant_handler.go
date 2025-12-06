package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type MerchantHandler struct {
	merchantService service.MerchantService
}

func NewMerchantHandler(merchantService service.MerchantService) *MerchantHandler {
	return &MerchantHandler{
		merchantService: merchantService,
	}
}

// Create godoc
// @Summary Create Merchant
// @Description Membuat merchant baru
// @Tags Merchants
// @Accept json
// @Produce json
// @Param request body dto.MerchantRequest true "Merchant data"
// @Success 201 {object} presenter.SuccessResponseSwagger{data=dto.MerchantResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /merchants [post]
func (h *MerchantHandler) Create(c *fiber.Ctx) error {
	var req dto.MerchantRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.merchantService.Create(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// GetAll godoc
// @Summary Get All Merchants
// @Description Mendapatkan daftar semua merchant
// @Tags Merchants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.MerchantListResponse}
// @Router /merchants [get]
func (h *MerchantHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	response := h.merchantService.GetAll(c.Context(), page, limit)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// GetById godoc
// @Summary Get Merchant by ID
// @Description Mendapatkan detail merchant berdasarkan ID
// @Tags Merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.MerchantResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Router /merchants/{id} [get]
func (h *MerchantHandler) GetById(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	response := h.merchantService.GetById(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// Update godoc
// @Summary Update Merchant
// @Description Update data merchant berdasarkan ID
// @Tags Merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Param request body dto.MerchantRequest true "Merchant data"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.MerchantResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /merchants/{id} [put]
func (h *MerchantHandler) Update(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}
	var req dto.MerchantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.merchantService.Update(c.Context(), c.Params("id"), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// Delete godoc
// @Summary Delete Merchant
// @Description Menghapus merchant berdasarkan ID
// @Tags Merchants
// @Accept json
// @Produce json
// @Param id path string true "Merchant ID"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /merchants/{id} [delete]
func (h *MerchantHandler) Delete(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}
	response := h.merchantService.Delete(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}
