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
)

type CustomerService interface {
	Create(ctx context.Context, req *dto.CustomerRequest) *presenter.Response
	GetAll(ctx context.Context, page, limit int) *presenter.Response
	GetById(ctx context.Context, id string) *presenter.Response
	Update(ctx context.Context, id string, req *dto.CustomerRequest) *presenter.Response
	Delete(ctx context.Context, id string) *presenter.Response
}

type customerService struct {
	customerRepo repositories.CustomerRepository
	cfg          *yaml.Config
}

func NewCustomerService(customerRepo repositories.CustomerRepository, cfg *yaml.Config) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
		cfg:          cfg,
	}
}

func (s *customerService) Create(ctx context.Context, req *dto.CustomerRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("customer_service_create", s.cfg.Logger.Enable)
	)

	if req.Email != "" {
		log.Info("checking if email already exists")
		existing, err := s.customerRepo.FindOneByEmail(ctx, req.Email)
		if err != nil {
			log.Error(fmt.Sprintf("error checking email: %v", err))
			return response.WithCode(500).WithError(errors.New("something went wrong"))
		}

		if existing != nil {
			log.Warn("email already exists")
			return response.WithCode(400).WithError(errors.New("email already exists"))
		}
	}

	// Helper function
	stringPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	customer := &entities.Customer{
		Name:  req.Name,
		Email: stringPtr(req.Email),
		Phone: stringPtr(req.Phone),
	}

	log.Info("creating customer")
	created, err := s.customerRepo.Create(ctx, customer)
	if err != nil {
		log.Error(fmt.Sprintf("error creating customer: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create customer"))
	}

	data := dto.ToCustomerResponse(created)
	return response.WithCode(201).WithData(data)
}

func (s *customerService) GetAll(ctx context.Context, page, limit int) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("customer_service_getall", s.cfg.Logger.Enable)
	)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	log.Info("fetching customers")
	customers, err := s.customerRepo.FindAll(ctx, limit, offset)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching customers: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to fetch customers"))
	}

	total, err := s.customerRepo.Count(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("error counting customers: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to count customers"))
	}

	data := dto.ToCustomerListResponse(customers, total, page, limit)
	return response.WithCode(200).WithData(data)
}

func (s *customerService) GetById(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("customer_service_getbyid", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching customer with id: %d", id))
	customer, err := s.customerRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching customer: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if customer == nil {
		log.Warn("customer not found")
		return response.WithCode(404).WithError(errors.New("customer not found"))
	}

	data := dto.ToCustomerResponse(customer)
	return response.WithCode(200).WithData(data)
}

func (s *customerService) Update(ctx context.Context, id string, req *dto.CustomerRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("customer_service_update", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching customer with id: %d", id))
	customer, err := s.customerRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching customer: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if customer == nil {
		log.Warn("customer not found")
		return response.WithCode(404).WithError(errors.New("customer not found"))
	}

	// Check email uniqueness if changed
	if req.Email != "" {
		currentEmail := ""
		if customer.Email != nil {
			currentEmail = *customer.Email
		}
		if req.Email != currentEmail {
			existing, err := s.customerRepo.FindOneByEmail(ctx, req.Email)
			if err != nil {
				log.Error(fmt.Sprintf("error checking email: %v", err))
				return response.WithCode(500).WithError(errors.New("something went wrong"))
			}
			if existing != nil {
				log.Warn("email already exists")
				return response.WithCode(400).WithError(errors.New("email already exists"))
			}
		}
	}

	customer.Name = req.Name
	customer.Email = stringPtr(req.Email)
	customer.Phone = stringPtr(req.Phone)

	log.Info("updating customer")
	updated, err := s.customerRepo.Update(ctx, customer)
	if err != nil {
		log.Error(fmt.Sprintf("error updating customer: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to update customer"))
	}

	data := dto.ToCustomerResponse(updated)
	return response.WithCode(200).WithData(data)
}

func (s *customerService) Delete(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("customer_service_delete", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("checking if customer exists: %d", id))
	customer, err := s.customerRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching customer: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if customer == nil {
		log.Warn("customer not found")
		return response.WithCode(404).WithError(errors.New("customer not found"))
	}

	log.Info("deleting customer")
	err = s.customerRepo.Delete(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error deleting customer: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to delete customer"))
	}

	return response.WithCode(200).WithData(map[string]string{"message": "customer deleted successfully"})
}
