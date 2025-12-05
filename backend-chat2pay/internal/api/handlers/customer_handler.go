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

func (h *CustomerHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	response := h.customerService.GetAll(c.Context(), page, limit)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

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
