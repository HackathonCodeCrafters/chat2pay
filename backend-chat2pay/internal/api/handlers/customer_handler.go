package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// Create godoc
// @Summary Create Customer
// @Description Membuat customer baru
// @Tags Customers
// @Accept json
// @Produce json
// @Param request body dto.CustomerRequest true "Customer data"
// @Success 201 {object} presenter.SuccessResponseSwagger{data=dto.CustomerResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /customers [post]
func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var req dto.CustomerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.customerService.Create(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// GetAll godoc
// @Summary Get All Customers
// @Description Mendapatkan daftar semua customer (Merchant Only)
// @Tags Customers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.CustomerListResponse}
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /customers [get]
func (h *CustomerHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	response := h.customerService.GetAll(c.Context(), page, limit)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// GetById godoc
// @Summary Get Customer by ID
// @Description Mendapatkan detail customer berdasarkan ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.CustomerResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetById(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	response := h.customerService.GetById(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// Update godoc
// @Summary Update Customer
// @Description Update data customer berdasarkan ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param request body dto.CustomerRequest true "Customer data"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.CustomerResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /customers/{id} [put]
func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	var req dto.CustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.customerService.Update(c.Context(), c.Params("id"), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// Delete godoc
// @Summary Delete Customer
// @Description Menghapus customer berdasarkan ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	response := h.customerService.Delete(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}
