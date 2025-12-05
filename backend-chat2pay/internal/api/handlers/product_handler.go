package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req dto.ProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.productService.Create(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if c.Query("merchant_id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	response := h.productService.GetAll(c.Context(), c.Query("merchant_id"), page, limit)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

func (h *ProductHandler) GetById(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	response := h.productService.GetById(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	var req dto.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.productService.Update(c.Context(), c.Params("id"), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	response := h.productService.Delete(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}
