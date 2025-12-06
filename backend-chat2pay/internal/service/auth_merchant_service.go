package service

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/dto"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/entities"
	"chat2pay/internal/middlewares/jwt"
	"chat2pay/internal/pkg/logger"
	"chat2pay/internal/repositories"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type MerchantAuthService interface {
	Register(ctx context.Context, req *dto.MerchantRegisterRequest) *presenter.Response
	Login(ctx context.Context, req *dto.MerchantLoginRequest) *presenter.Response
}

type merchantAuthService struct {
	merchantUserRepo repositories.MerchantUserRepository
	merchantRepo     repositories.MerchantRepository
	authMdwr         jwt.AuthMiddleware
	cfg              *yaml.Config
}

func NewMerchantAuthService(
	merchantUserRepo repositories.MerchantUserRepository,
	merchantRepo repositories.MerchantRepository,
	authMdwr jwt.AuthMiddleware,
	cfg *yaml.Config,
) MerchantAuthService {
	return &merchantAuthService{
		merchantUserRepo: merchantUserRepo,
		merchantRepo:     merchantRepo,
		authMdwr:         authMdwr,
		cfg:              cfg,
	}
}

func (s *merchantAuthService) Register(ctx context.Context, req *dto.MerchantRegisterRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("merchant_auth_register", s.cfg.Logger.Enable)
	)

	log.Info("checking if email already exists")
	existingUser, err := s.merchantUserRepo.FindOneByEmail(ctx, req.Email)
	if err != nil {
		log.Error(fmt.Sprintf("error checking email: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if existingUser != nil {
		log.Warn("email already registered")
		return response.WithCode(400).WithError(errors.New("email already registered"))
	}

	log.Info("hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(fmt.Sprintf("error hashing password: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to hash password"))
	}

	log.Info("creating merchant")
	merchant := &entities.Merchant{
		Name:      req.MerchantName,
		LegalName: stringPtr(req.LegalName),
		Email:     req.Email,
		Phone:     stringPtr(req.Phone),
	}

	createdMerchant, err := s.merchantRepo.Create(ctx, merchant)
	if err != nil {
		log.Error(fmt.Sprintf("error creating merchant: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create merchant"))
	}

	log.Info("creating merchant user")
	merchantUser := &entities.MerchantUser{
		MerchantID:   createdMerchant.ID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "owner",
		IsActive:     true,
	}

	createdUser, err := s.merchantUserRepo.Create(ctx, merchantUser)
	if err != nil {
		log.Error(fmt.Sprintf("error creating merchant user: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create merchant user"))
	}

	log.Info("generating token")
	token, err := s.authMdwr.GenerateToken(createdUser.ID.String(), createdUser.Email, "merchant")
	if err != nil {
		log.Error(fmt.Sprintf("error generating token: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to generate token"))
	}

	createdUser.Merchant = createdMerchant
	data := dto.ToMerchantAuthResponse(createdUser, *token)
	return response.WithCode(201).WithData(data)
}

func (s *merchantAuthService) Login(ctx context.Context, req *dto.MerchantLoginRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("merchant_auth_login", s.cfg.Logger.Enable)
	)

	log.Info("finding user by email")
	merchantUser, err := s.merchantUserRepo.FindOneByEmail(ctx, req.Email)
	if err != nil {
		log.Error(fmt.Sprintf("error finding user: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if merchantUser == nil {
		log.Warn("user not found")
		return response.WithCode(401).WithError(errors.New("invalid email or password"))
	}

	if !merchantUser.IsActive {
		log.Warn("user is not active")
		return response.WithCode(401).WithError(errors.New("account is not active"))
	}

	log.Info("comparing password")
	err = bcrypt.CompareHashAndPassword([]byte(merchantUser.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Warn("invalid password")
		return response.WithCode(401).WithError(errors.New("invalid email or password"))
	}

	log.Info("generating token")
	token, err := s.authMdwr.GenerateToken(merchantUser.ID.String(), merchantUser.Email, "merchant")
	if err != nil {
		log.Error(fmt.Sprintf("error generating token: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to generate token"))
	}

	data := dto.ToMerchantAuthResponse(merchantUser, *token)
	return response.WithCode(200).WithData(data)
}
