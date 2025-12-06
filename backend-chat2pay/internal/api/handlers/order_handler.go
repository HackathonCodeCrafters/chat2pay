package handlers

import (
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Create godoc
// @Summary Create Order
// @Description Membuat order baru
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.OrderRequest true "Order data"
// @Success 201 {object} presenter.SuccessResponseSwagger{data=dto.OrderResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /orders [post]
func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var req dto.OrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.orderService.Create(c.Context(), &req)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// GetAll godoc
// @Summary Get All Orders
// @Description Mendapatkan daftar semua order
// @Tags Orders
// @Accept json
// @Produce json
// @Param merchant_id query string false "Filter by Merchant ID"
// @Param customer_id query string false "Filter by Customer ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.OrderListResponse}
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /orders [get]
func (h *OrderHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	merchantId := c.Query("merchant_id")
	customerId := c.Query("customer_id")

	response := h.orderService.GetAll(c.Context(), merchantId, customerId, page, limit)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// GetById godoc
// @Summary Get Order by ID
// @Description Mendapatkan detail order berdasarkan ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.OrderResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /orders/{id} [get]
func (h *OrderHandler) GetById(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}
	response := h.orderService.GetById(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// UpdateStatus godoc
// @Summary Update Order Status
// @Description Update status order (Merchant Only)
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.OrderStatusUpdateRequest true "Status update"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.OrderResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}
	var req dto.OrderStatusUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.orderService.UpdateStatus(c.Context(), c.Params("id"), req.Status)

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}

// Delete godoc
// @Summary Delete Order
// @Description Menghapus order (Merchant Only)
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /orders/{id} [delete]
func (h *OrderHandler) Delete(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return c.Status(400).JSON(presenter.ErrorResponse(fiber.ErrBadRequest))
	}

	response := h.orderService.Delete(c.Context(), c.Params("id"))

	if response.Errors != nil {
		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}
