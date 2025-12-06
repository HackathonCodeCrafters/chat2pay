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

// Create godoc
// @Summary Create Product
// @Description Membuat produk baru
// @Tags Products
// @Accept json
// @Produce json
// @Param request body dto.ProductRequest true "Product data"
// @Success 201 {object} presenter.SuccessResponseSwagger{data=dto.ProductResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /products [post]
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

// GetAll godoc
// @Summary Get All Products
// @Description Mendapatkan daftar produk berdasarkan merchant
// @Tags Products
// @Accept json
// @Produce json
// @Param merchant_id query string true "Merchant ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.ProductListResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Router /products [get]
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

// GetById godoc
// @Summary Get Product by ID
// @Description Mendapatkan detail produk berdasarkan ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.ProductResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Router /products/{id} [get]
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

// Update godoc
// @Summary Update Product
// @Description Update data produk berdasarkan ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dto.ProductRequest true "Product data"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.ProductResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /products/{id} [put]
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

// Delete godoc
// @Summary Delete Product
// @Description Menghapus produk berdasarkan ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} presenter.SuccessResponseSwagger
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Failure 401 {object} presenter.ErrorResponseSwagger
// @Failure 403 {object} presenter.ErrorResponseSwagger
// @Failure 404 {object} presenter.ErrorResponseSwagger
// @Security BearerAuth
// @Router /products/{id} [delete]
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

// AskProduct godoc
// @Summary Ask Product (AI Search)
// @Description Mencari produk menggunakan AI/LLM dengan natural language query
// @Tags AI-LLM
// @Accept json
// @Produce json
// @Param request body dto.AskProduct true "Prompt pencarian"
// @Success 200 {object} presenter.SuccessResponseSwagger{data=dto.LLMResponse}
// @Failure 400 {object} presenter.ErrorResponseSwagger
// @Router /products/ask [post]
func (h *ProductHandler) AskProduct(c *fiber.Ctx) error {

	var req dto.AskProduct

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(presenter.ErrorResponse(err))
	}

	response := h.productService.AskProduct(c.Context(), &req)

	if response.Errors != nil {

		return c.Status(response.Code).JSON(presenter.ErrorResponse(response.Errors))
	}

	return c.Status(response.Code).JSON(presenter.SuccessResponse(response.Data))
}
