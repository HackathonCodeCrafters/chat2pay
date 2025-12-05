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

type MerchantService interface {
	Create(ctx context.Context, req *dto.MerchantRequest) *presenter.Response
	GetAll(ctx context.Context, page, limit int) *presenter.Response
	GetById(ctx context.Context, id string) *presenter.Response
	Update(ctx context.Context, id string, req *dto.MerchantRequest) *presenter.Response
	Delete(ctx context.Context, id string) *presenter.Response
}

type merchantService struct {
	merchantRepo repositories.MerchantRepository
	cfg          *yaml.Config
}

func NewMerchantService(merchantRepo repositories.MerchantRepository, cfg *yaml.Config) MerchantService {
	return &merchantService{
		merchantRepo: merchantRepo,
		cfg:          cfg,
	}
}

func (s *merchantService) Create(ctx context.Context, req *dto.MerchantRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("merchant_service_create", s.cfg.Logger.Enable)
	)

	log.Info("checking if email already exists")
	existing, err := s.merchantRepo.FindOneByEmail(ctx, req.Email)
	if err != nil {
		log.Error(fmt.Sprintf("error checking email: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if existing != nil {
		log.Warn("email already exists")
		return response.WithCode(400).WithError(errors.New("email already exists"))
	}

	merchant := &entities.Merchant{
		Name:      req.Name,
		LegalName: stringPtr(req.LegalName),
		Email:     req.Email,
		Phone:     stringPtr(req.Phone),
		Status:    "pending_verification",
	}

	if req.Status != "" {
		merchant.Status = req.Status
	}

	log.Info("creating merchant")
	created, err := s.merchantRepo.Create(ctx, merchant)
	if err != nil {
		log.Error(fmt.Sprintf("error creating merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create merchant"))
	}

	data := dto.ToMerchantResponse(created)
	return response.WithCode(201).WithData(data)
}

func (s *merchantService) GetAll(ctx context.Context, page, limit int) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("merchant_service_getall", s.cfg.Logger.Enable)
	)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	log.Info("fetching merchants")
	merchants, err := s.merchantRepo.FindAll(ctx, limit, offset)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching merchants: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to fetch merchants"))
	}

	total, err := s.merchantRepo.Count(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("error counting merchants: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to count merchants"))
	}

	data := dto.ToMerchantListResponse(merchants, total, page, limit)
	return response.WithCode(200).WithData(data)
}

func (s *merchantService) GetById(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("merchant_service_getbyid", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching merchant with id: %d", id))
	merchant, err := s.merchantRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if merchant == nil {
		log.Warn("merchant not found")
		return response.WithCode(404).WithError(errors.New("merchant not found"))
	}

	data := dto.ToMerchantResponse(merchant)
	return response.WithCode(200).WithData(data)
}

func (s *merchantService) Update(ctx context.Context, id string, req *dto.MerchantRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("merchant_service_update", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("fetching merchant with id: %d", id))
	merchant, err := s.merchantRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if merchant == nil {
		log.Warn("merchant not found")
		return response.WithCode(404).WithError(errors.New("merchant not found"))
	}

	if req.Email != merchant.Email {
		existing, err := s.merchantRepo.FindOneByEmail(ctx, req.Email)
		if err != nil {
			log.Error(fmt.Sprintf("error checking email: %v", err))
			return response.WithCode(500).WithError(errors.New("something went wrong"))
		}
		if existing != nil {
			log.Warn("email already exists")
			return response.WithCode(400).WithError(errors.New("email already exists"))
		}
	}

	merchant.Name = req.Name
	merchant.LegalName = stringPtr(req.LegalName)
	merchant.Email = req.Email
	merchant.Phone = stringPtr(req.Phone)

	if req.Status != "" {
		merchant.Status = req.Status
	}

	log.Info("updating merchant")
	updated, err := s.merchantRepo.Update(ctx, merchant)
	if err != nil {
		log.Error(fmt.Sprintf("error updating merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to update merchant"))
	}

	data := dto.ToMerchantResponse(updated)
	return response.WithCode(200).WithData(data)
}

func (s *merchantService) Delete(ctx context.Context, id string) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("merchant_service_delete", s.cfg.Logger.Enable)
	)

	log.Info(fmt.Sprintf("checking if merchant exists: %d", id))
	merchant, err := s.merchantRepo.FindOneById(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error fetching merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if merchant == nil {
		log.Warn("merchant not found")
		return response.WithCode(404).WithError(errors.New("merchant not found"))
	}

	log.Info("deleting merchant")
	err = s.merchantRepo.Delete(ctx, id)
	if err != nil {
		log.Error(fmt.Sprintf("error deleting merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to delete merchant"))
	}

	return response.WithCode(200).WithData(map[string]string{"message": "merchant deleted successfully"})
}
