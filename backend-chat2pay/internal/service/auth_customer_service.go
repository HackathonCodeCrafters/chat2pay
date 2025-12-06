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

type CustomerAuthService interface {
	Register(ctx context.Context, req *dto.CustomerRegisterRequest) *presenter.Response
	Login(ctx context.Context, req *dto.CustomerLoginRequest) *presenter.Response
}

type customerAuthService struct {
	customerRepo repositories.CustomerRepository
	authMdwr     jwt.AuthMiddleware
	cfg          *yaml.Config
}

func NewCustomerAuthService(
	customerRepo repositories.CustomerRepository,
	authMdwr jwt.AuthMiddleware,
	cfg *yaml.Config,
) CustomerAuthService {
	return &customerAuthService{
		customerRepo: customerRepo,
		authMdwr:     authMdwr,
		cfg:          cfg,
	}
}

func (s *customerAuthService) Register(ctx context.Context, req *dto.CustomerRegisterRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("customer_auth_register", s.cfg.Logger.Enable)
	)

	log.Info("checking if email already exists")
	existingCustomer, err := s.customerRepo.FindOneByEmail(ctx, req.Email)
	if err != nil {
		log.Error(fmt.Sprintf("error checking email: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if existingCustomer != nil {
		log.Warn("email already registered")
		return response.WithCode(400).WithError(errors.New("email already registered"))
	}

	log.Info("hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(fmt.Sprintf("error hashing password: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to hash password"))
	}

	hashedPasswordStr := string(hashedPassword)
	log.Info("creating customer")
	customer := &entities.Customer{
		Name:         req.Name,
		Email:        &req.Email,
		Phone:        &req.Phone,
		PasswordHash: &hashedPasswordStr,
	}

	createdCustomer, err := s.customerRepo.Create(ctx, customer)
	if err != nil {
		log.Error(fmt.Sprintf("error creating customer: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to create customer"))
	}

	log.Info("generating token")
	token, err := s.authMdwr.GenerateToken(createdCustomer.ID.String(), *createdCustomer.Email, "customer")
	if err != nil {
		log.Error(fmt.Sprintf("error generating token: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to generate token"))
	}

	data := dto.ToCustomerAuthResponse(createdCustomer, *token)
	return response.WithCode(201).WithData(data)
}

func (s *customerAuthService) Login(ctx context.Context, req *dto.CustomerLoginRequest) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("customer_auth_login", s.cfg.Logger.Enable)
	)

	log.Info("finding customer by email")
	customer, err := s.customerRepo.FindOneByEmail(ctx, req.Email)
	if err != nil {
		log.Error(fmt.Sprintf("error finding customer: %v", err))
		return response.WithCode(500).WithError(errors.New("something went wrong"))
	}

	if customer == nil || customer.PasswordHash == nil {
		log.Warn("customer not found or no password set")
		return response.WithCode(401).WithError(errors.New("invalid email or password"))
	}

	log.Info("comparing password")
	err = bcrypt.CompareHashAndPassword([]byte(*customer.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Warn("invalid password")
		return response.WithCode(401).WithError(errors.New("invalid email or password"))
	}

	log.Info("generating token")
	token, err := s.authMdwr.GenerateToken(customer.ID.String(), *customer.Email, "customer")
	if err != nil {
		log.Error(fmt.Sprintf("error generating token: %v", err))
		return response.WithCode(500).WithError(errors.New("failed to generate token"))
	}

	data := dto.ToCustomerAuthResponse(customer, *token)
	return response.WithCode(200).WithData(data)
}
