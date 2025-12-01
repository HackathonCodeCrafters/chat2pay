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
	"fmt"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, user *entities.User) *presenter.Response
	GetProfile(ctx context.Context, userId string) *presenter.Response
}

type authService struct {
	jwtMdwr  jwt.AuthMiddleware
	userRepo repositories.UserRepository
	cfg      *yaml.Config
}

func NewAuthService(jwtMdwr jwt.AuthMiddleware, userRepo repositories.UserRepository, cfg *yaml.Config) AuthService {
	return &authService{
		jwtMdwr:  jwtMdwr,
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Login(ctx context.Context, user *entities.User) *presenter.Response {
	var (
		response = presenter.Response{}
		log      = logger.NewLog("login_service", s.cfg.Logger.Enable)
	)

	log.Info("fetching user data from db")

	userData, err := s.userRepo.FindOneByEmail(ctx, user.Email)

	if err != nil {
		log.Error(fmt.Sprintf(`error fetching user data to db got %s`, err))
		return response.WithCode(500).WithError(errors.New("something went wrong!"))
	}

	if userData == nil {
		log.Warn("user not found in db")
		return response.WithCode(401).WithError(errors.New("email or password is wrong"))
	}

	log.Info("comparing user password")

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
	if err != nil {
		log.Error(fmt.Sprintf(`error comparing user password got %s`, err))
		return response.WithCode(401).WithError(errors.New("email or password is wrong"))
	}

	log.Info("generating jwt token")
	token, err := s.jwtMdwr.GenerateToken(userData)

	if err != nil {
		log.Error(fmt.Sprintf(`error generating jwt token got %s`, err))
		return response.WithCode(500).WithError(errors.New("something went wrong!"))
	}

	data := dto.AuthSucces(userData, *token)

	return response.WithCode(200).WithData(data)
}
func (s *authService) GetProfile(ctx context.Context, userId string) *presenter.Response {
	var (
		log      = logger.NewLog("login_handler", s.cfg.Logger.Enable)
		response = presenter.Response{}
	)

	log.Info("fetching user data from db")

	userData, err := s.userRepo.FindOneById(ctx, userId)

	if err != nil {
		log.Error(fmt.Sprintf(`error fetching user data to db got %s`, err))
		return response.WithCode(500).WithError(errors.New("something went wrong!"))
	}

	if userData == nil {
		log.Warn("user not found in db")
		return response.WithCode(401).WithError(errors.New("user not found"))
	}

	data := dto.GetProfileSuccess(userData)

	return response.WithCode(200).WithData(data)
}
