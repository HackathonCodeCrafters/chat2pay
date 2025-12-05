package service

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/entities"
	"chat2pay/internal/pkg/logger"
	"chat2pay/internal/repositories"
	"context"
	"errors"
	"fmt"
	"time"
)

type OrderService interface {
	Create(ctx context.Context, req *dto.OrderRequest) *presenter.Response
	GetAll(ctx context.Context, merchantId, customerId string, page, limit int) *presenter.Response
	GetById(ctx context.Context, id string) *presenter.Response
	UpdateStatus(ctx context.Context, id string, status string) *presenter.Response
	Delete(ctx context.Context, id string) *presenter.Response
}

type orderService struct {
	orderRepo    repositories.OrderRepository
	productRepo  repositories.ProductRepository
	customerRepo repositories.CustomerRepository
	merchantRepo repositories.MerchantRepository
	cfg          *yaml.Config
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	productRepo repositories.ProductRepository,
	customerRepo repositories.CustomerRepository,
	merchantRepo repositories.MerchantRepository,
	cfg *yaml.Config,
) OrderService {
	return &orderService{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
		merchantRepo: merchantRepo,
		cfg:          cfg,
	}
}

func (s *orderService) Create(ctx context.Context, req *dto.OrderRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("order_service_create", s.cfg.Logger.Enable)
	)

	log.Info("validating customer")
	customer, err := s.customerRepo.FindOneById(ctx, req.CustomerID)
	if err != nil {
		log.Error(fmt.Sprintf("error checking customer: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}
	if customer == nil {
		log.Warn("customer not found")
		return response.WithCode(404).WithError(errors.New("customer not found"))
	}

	log.Info("validating merchant")
	merchant, err := s.merchantRepo.FindOneById(ctx, req.MerchantID)
	if err != nil {
		log.Error(fmt.Sprintf("error checking merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}
	if merchant == nil {
		log.Warn("merchant not found")
		return response.WithCode(404).WithError(errors.New("merchant not found"))
	}

	var orderItems []entities.OrderItem
	var subtotal float64 = 0

	log.Info("validating products and calculating total")
	for _, item := range req.Items {
		product, err := s.productRepo.FindOneById(ctx, item.ProductID)
		if err != nil {
			log.Error(fmt.Sprintf("error fetching product: %v", err))
			return response.WithCode(500).WithError(errors.New("something went wrong"))
		}
		if product == nil {
			log.Warn(fmt.Sprintf("product with id %d not found", item.ProductID))
			return response.WithCode(404).WithError(fmt.Errorf("product with id %d not found", item.ProductID))
		}

		if product.Stock < item.Quantity {
			log.Warn(fmt.Sprintf("insufficient stock for product %s", product.Name))
			return response.WithCode(400).WithError(fmt.Errorf("insufficient stock for product: %s", product.Name))
		}

		itemTotal := product.Price * float64(item.Quantity)
		subtotal += itemTotal

		orderItems = append(orderItems, entities.OrderItem{
			ProductID:           item.ProductID,
			ProductNameSnapshot: product.Name,
			UnitPrice:           product.Price,
			Qty:                 item.Quantity,
			TotalPrice:          itemTotal,
		})
	}

	totalAmount := subtotal + req.ShippingAmount - req.DiscountAmount

	orderNumber := fmt.Sprintf("ORD-%d-%d", req.MerchantID, time.Now().Unix())

	order := &entities.Order{
		OrderNumber:    orderNumber,
		CustomerID:     req.CustomerID,
		MerchantID:     req.MerchantID,
		OutletID:       req.OutletID,
		Status:         "pending",
		SubtotalAmount: subtotal,
		ShippingAmount: req.ShippingAmount,
		DiscountAmount: req.DiscountAmount,
		TotalAmount:    totalAmount,
	}

	log.Info("creating order with items")
	created, err := s.orderRepo.CreateWithItems(ctx, order, orderItems)
	if err != nil {
		log.Error(fmt.Sprintf("error creating order: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create order"))
	}

	log.Info("updating product stocks")
	for _, item := range req.Items {
		err := s.productRepo.UpdateStock(ctx, item.ProductID, -item.Quantity)
		if err != nil {
			log.Error(fmt.Sprintf("error updating stock for product %d: %v", item.ProductID, err))
		}
	}

	createdOrder, _ := s.orderRepo.FindOneById(ctx, created.ID)
	data := dto.ToOrderResponse(createdOrder)
	return response.WithCode(201).WithData(data)
}

func (s *orderService) GetAll(ctx context.Context, merchantId, customerId string, page, limit int) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("order_service_getall", s.cfg.Logger.Enable)
	)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	log.Info("fetching orders")
	orders, err := s.orderRepo.FindAll(ctx, merchantId, customerId, limit, offset)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching orders: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to fetch orders"))
	}

	total, err := s.orderRepo.Count(ctx, merchantId, customerId)
	if err != nil {
		log.Error(fmt.Sprintf("error counting orders: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to count orders"))
	}

	data := dto.ToOrderListResponse(orders, total, page, limit)
	return response.WithCode(200).WithData(data)
}

func (s *orderService) GetById(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("order_service_getbyid", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching order with id: %d", id))
	order, err := s.orderRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching order: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if order == nil {
		log.Warn("order not found")
		return response.WithCode(404).WithError(errors.New("order not found"))
	}

	data := dto.ToOrderResponse(order)
	return response.WithCode(200).WithData(data)
}

func (s *orderService) UpdateStatus(ctx context.Context, id string, status string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("order_service_updatestatus", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching order with id: %d", id))
	order, err := s.orderRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching order: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if order == nil {
		log.Warn("order not found")
		return response.WithCode(404).WithError(errors.New("order not found"))
	}

	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		log.Warn("invalid status")
		return response.WithCode(400).WithError(errors.New("invalid status"))
	}

	log.Info("updating order status")
	err = s.orderRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		log.Error(fmt.Sprintf("error updating order status: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to update order status"))
	}

	updatedOrder, _ := s.orderRepo.FindOneById(ctx, id)
	data := dto.ToOrderResponse(updatedOrder)
	return response.WithCode(200).WithData(data)
}

func (s *orderService) Delete(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("order_service_delete", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("checking if order exists: %d", id))
	order, err := s.orderRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching order: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if order == nil {
		log.Warn("order not found")
		return response.WithCode(404).WithError(errors.New("order not found"))
	}

	log.Info("deleting order")
	err = s.orderRepo.Delete(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error deleting order: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to delete order"))
	}

	return response.WithCode(200).WithData(map[string]string{"message": "order deleted successfully"})
}
