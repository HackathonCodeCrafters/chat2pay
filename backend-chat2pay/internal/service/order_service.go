package service

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/entities"
	"chat2pay/internal/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID string, req *dto.CreateOrderRequest) *presenter.Response
	GetOrderByID(ctx context.Context, id string) *presenter.Response
	GetCustomerOrders(ctx context.Context, customerID string, page, limit int) *presenter.Response
	GetMerchantOrders(ctx context.Context, merchantID string, page, limit int) *presenter.Response
	UpdateOrderStatus(ctx context.Context, orderID, status string) *presenter.Response
	UpdateTrackingNumber(ctx context.Context, orderID, trackingNumber string) *presenter.Response
}

type orderService struct {
	cfg         *yaml.Config
	orderRepo   repositories.OrderRepository
	productRepo repositories.ProductRepository
}

func NewOrderService(cfg *yaml.Config, orderRepo repositories.OrderRepository, productRepo repositories.ProductRepository) OrderService {
	return &orderService{
		cfg:         cfg,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, customerID string, req *dto.CreateOrderRequest) *presenter.Response {
	response := presenter.NewResponse()

	// Calculate subtotal and validate products
	var subtotal float64
	var merchantID string
	orderItems := make([]entities.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		product, err := s.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return response.WithCode(404).WithError(errors.New("product not found: " + item.ProductID))
		}

		if product.Stock < item.Quantity {
			return response.WithCode(400).WithError(errors.New("insufficient stock for: " + product.Name))
		}

		// All products must be from the same merchant
		if merchantID == "" {
			merchantID = product.MerchantID
		} else if merchantID != product.MerchantID {
			return response.WithCode(400).WithError(errors.New("all products must be from the same merchant"))
		}

		itemSubtotal := product.Price * float64(item.Quantity)
		subtotal += itemSubtotal

		orderItems = append(orderItems, entities.OrderItem{
			ID:           uuid.New().String(),
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     item.Quantity,
			Subtotal:     itemSubtotal,
		})
	}

	total := subtotal + req.ShippingCost

	// Create order
	order := &entities.Order{
		ID:                 uuid.New().String(),
		CustomerID:         customerID,
		MerchantID:         merchantID,
		Status:             entities.OrderStatusPending,
		Subtotal:           subtotal,
		ShippingCost:       req.ShippingCost,
		Total:              total,
		Courier:            &req.Courier,
		CourierService:     &req.CourierService,
		ShippingEtd:        &req.ShippingEtd,
		ShippingAddress:    &req.ShippingAddress,
		ShippingCity:       &req.ShippingCity,
		ShippingProvince:   &req.ShippingProvince,
		ShippingPostalCode: &req.ShippingPostalCode,
		PaymentStatus:      entities.PaymentStatusPending,
	}

	if req.Notes != "" {
		order.Notes = &req.Notes
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return response.WithCode(500).WithError(errors.New("failed to create order"))
	}

	// Create order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := s.orderRepo.CreateItem(ctx, &orderItems[i]); err != nil {
			return response.WithCode(500).WithError(errors.New("failed to create order item"))
		}
	}

	// Update product stock
	for _, item := range req.Items {
		product, _ := s.productRepo.FindByID(ctx, item.ProductID)
		newStock := product.Stock - item.Quantity
		s.productRepo.UpdateStock(ctx, item.ProductID, newStock)
	}

	order.Items = orderItems
	return response.WithCode(201).WithData(dto.ToOrderResponse(order))
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) *presenter.Response {
	response := presenter.NewResponse()

	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return response.WithCode(404).WithError(errors.New("order not found"))
	}

	items, _ := s.orderRepo.GetOrderItems(ctx, id)
	order.Items = items

	return response.WithCode(200).WithData(dto.ToOrderResponse(order))
}

func (s *orderService) GetCustomerOrders(ctx context.Context, customerID string, page, limit int) *presenter.Response {
	response := presenter.NewResponse()

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, err := s.orderRepo.FindByCustomerID(ctx, customerID, page, limit)
	if err != nil {
		return response.WithCode(500).WithError(errors.New("failed to get orders"))
	}

	// Load items for each order
	for i := range orders {
		items, _ := s.orderRepo.GetOrderItems(ctx, orders[i].ID)
		orders[i].Items = items
	}

	total, _ := s.orderRepo.CountByCustomerID(ctx, customerID)
	return response.WithCode(200).WithData(dto.ToOrderListResponse(orders, total, page, limit))
}

func (s *orderService) GetMerchantOrders(ctx context.Context, merchantID string, page, limit int) *presenter.Response {
	response := presenter.NewResponse()

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, err := s.orderRepo.FindByMerchantID(ctx, merchantID, page, limit)
	if err != nil {
		return response.WithCode(500).WithError(errors.New("failed to get orders"))
	}

	// Load items for each order
	for i := range orders {
		items, _ := s.orderRepo.GetOrderItems(ctx, orders[i].ID)
		orders[i].Items = items
	}

	total, _ := s.orderRepo.CountByMerchantID(ctx, merchantID)
	return response.WithCode(200).WithData(dto.ToOrderListResponse(orders, total, page, limit))
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID, status string) *presenter.Response {
	response := presenter.NewResponse()

	err := s.orderRepo.UpdateStatus(ctx, orderID, status)
	if err != nil {
		return response.WithCode(500).WithError(errors.New("failed to update order status"))
	}

	return response.WithCode(200).WithData(map[string]string{"message": "Order status updated"})
}

func (s *orderService) UpdateTrackingNumber(ctx context.Context, orderID, trackingNumber string) *presenter.Response {
	response := presenter.NewResponse()

	err := s.orderRepo.UpdateTrackingNumber(ctx, orderID, trackingNumber)
	if err != nil {
		return response.WithCode(500).WithError(errors.New("failed to update tracking number"))
	}

	return response.WithCode(200).WithData(map[string]string{"message": "Tracking number updated"})
}
