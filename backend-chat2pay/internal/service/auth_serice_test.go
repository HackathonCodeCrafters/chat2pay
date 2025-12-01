package service

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/entities"
	"chat2pay/internal/middlewares/jwt"
	"chat2pay/internal/repositories"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var (
	cfg            = &yaml.Config{}
	userRepository = repositories.UserRepositoryMock{Mock: mock.Mock{}}
	//mdwr           = jwt.JWTMock{Mock: mock.Mock{}}
	jwtMiddleware = jwt.NewAuthMiddleware(&userRepository, cfg)
	loginService  = NewAuthService(jwtMiddleware, &userRepository, cfg)
)

func TestMain(m *testing.M) {
	cfg, _ = yaml.NewConfig()

	m.Run()
}

func TestAuthService_LoginFailed(t *testing.T) {
	t.Run("Get user got error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		userRepository.Mock.On("FindOneByEmail", "inierror@gmail.com").Return("error banh")

		login := loginService.Login(ctx, &entities.User{Email: "inierror@gmail.com", Password: "icikiwir"})
		assert.Equal(t, 500, login.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		userRepository.Mock.On("FindOneByEmail", "inigakada@gmail.com").Return(nil)

		login := loginService.Login(ctx, &entities.User{Email: "inigakada@gmail.com", Password: "icikiwir"})
		assert.Equal(t, 401, login.Code)
	})

	t.Run("Error comparing password", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		userRepository.Mock.On("FindOneByEmail", "iniada@gmail.com").Return(entities.User{Email: "iniada@gmail.com", Password: "$2a$10$sTkXJOBo2n3lWttFtsmnTebwkgzOhr9oisfy9H7EES0PCqalMgASm"})

		login := loginService.Login(ctx, &entities.User{Email: "iniada@gmail.com", Password: "aselole"})
		assert.Equal(t, 401, login.Code)
	})

	t.Run("Error gen password", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		userRepository.Mock.On("FindOneByEmail", "iniada@gmail.com").Return(entities.User{Email: "iniada@gmail.com", Password: "$2a$10$sTkXJOBo2n3lWttFtsmnTebwkgzOhr9oisfy9H7EES0PCqalMgASm"})

		login := loginService.Login(ctx, &entities.User{Email: "iniada@gmail.com", Password: "aselole"})
		assert.Equal(t, 401, login.Code)
	})

	//t.Run("Error generating jwt", func(t *testing.T) {
	//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//	defer cancel()
	//
	//	user := entities.User{Email: "iniada@gmail.com", Password: "$2a$10$sTkXJOBo2n3lWttFtsmnTebwkgzOhr9oisfy9H7EES0PCqalMgASm"}
	//
	//	mdwr.Mock.On("GenerateToken", user).Return("ini token jwt")
	//	userRepository.Mock.On("FindOneByEmail", "iniada@gmail.com").Return(user)
	//
	//	login := loginService.Login(ctx, &entities.User{Email: "iniada@gmail.com", Password: "aselole"})
	//	assert.Equal(t, 500, login.Code)
	//})
}

func TestAuthService_LoginSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userRepository.Mock.On("FindOneByEmail", "iniada@gmail.com").Return(entities.User{Email: "iniada@gmail.com", Password: "$2a$10$sTkXJOBo2n3lWttFtsmnTebwkgzOhr9oisfy9H7EES0PCqalMgASm"})

	login := loginService.Login(ctx, &entities.User{Email: "iniada@gmail.com", Password: "icikiwir"})
	assert.Equal(t, 200, login.Code)
}

func TestAuthService_GetProfileFailed(t *testing.T) {
	t.Run("Get user got error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		userRepository.Mock.On("FindOneById", "ini_bukan_id").Return("error banh")

		profile := loginService.GetProfile(ctx, "ini_bukan_id")
		assert.Equal(t, 500, profile.Code)
	})

	t.Run("User not found", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		userRepository.Mock.On("FindOneById", "ini_nggak_ada").Return(nil)

		profile := loginService.GetProfile(ctx, "ini_nggak_ada")
		assert.Equal(t, 401, profile.Code)
	})
}

func TestAuthService_GetProfileSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userRepository.Mock.On("FindOneById", "ini_id").Return(entities.User{
		ID:           "ini_id",
		Email:        "iniada@gmail.com",
		FullName:     "icikiwir",
		Age:          99,
		MobileNumber: "08123456789",
		Password:     "$2a$10$sTkXJOBo2n3lWttFtsmnTebwkgzOhr9oisfy9H7EES0PCqalMgASm",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		DeletedAt:    nil,
	})

	login := loginService.GetProfile(ctx, "ini_id")
	assert.Equal(t, 200, login.Code)
}
