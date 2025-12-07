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
	return &OrderHandler{orderService: orderService}
}

// CreateOrder godoc
// @Summary Create new order
// @Description Create a new order (checkout)
// @Tags Orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreateOrderRequest true "Order data"
// @Success 201 {object} presenter.SuccessResponseSwagger
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	customerIDVal := c.Locals("customer_id")
	if customerIDVal == nil {
		return c.Status(401).JSON(presenter.ErrorResponse(fiber.NewError(401, "Unauthorized: customer_id not found")))
	}
	customerID := customerIDVal.(string)

	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	result := h.orderService.CreateOrder(c.Context(), customerID, &req)
	return c.Status(result.Code).JSON(result)
}

// GetOrder godoc
// @Summary Get order by ID
// @Description Get order details
// @Tags Orders
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	result := h.orderService.GetOrderByID(c.Context(), id)
	return c.Status(result.Code).JSON(result)
}

// GetCustomerOrders godoc
// @Summary Get customer orders
// @Description Get all orders for current customer
// @Tags Orders
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /orders/customer [get]
func (h *OrderHandler) GetCustomerOrders(c *fiber.Ctx) error {
	customerIDVal := c.Locals("customer_id")
	if customerIDVal == nil {
		return c.Status(401).JSON(presenter.ErrorResponse(fiber.NewError(401, "Unauthorized: customer_id not found")))
	}
	customerID := customerIDVal.(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result := h.orderService.GetCustomerOrders(c.Context(), customerID, page, limit)
	return c.Status(result.Code).JSON(result)
}

// GetMerchantOrders godoc
// @Summary Get merchant orders
// @Description Get all orders for current merchant
// @Tags Orders
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /orders/merchant [get]
func (h *OrderHandler) GetMerchantOrders(c *fiber.Ctx) error {
	merchantIDVal := c.Locals("merchant_id")
	if merchantIDVal == nil {
		return c.Status(401).JSON(presenter.ErrorResponse(fiber.NewError(401, "Unauthorized: merchant_id not found")))
	}
	merchantID := merchantIDVal.(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result := h.orderService.GetMerchantOrders(c.Context(), merchantID, page, limit)
	return c.Status(result.Code).JSON(result)
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update order status (merchant only)
// @Tags Orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "Status data"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Router /orders/{id}/status [patch]
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	if req.TrackingNumber != "" {
		result := h.orderService.UpdateTrackingNumber(c.Context(), id, req.TrackingNumber)
		return c.Status(result.Code).JSON(result)
	}

	result := h.orderService.UpdateOrderStatus(c.Context(), id, req.Status)
	return c.Status(result.Code).JSON(result)
}
