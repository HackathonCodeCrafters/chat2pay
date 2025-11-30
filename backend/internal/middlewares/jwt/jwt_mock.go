package jwt

import (
	"chat2pay/internal/entities"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

type JWTMock struct {
	Mock mock.Mock
}

func (m *JWTMock) ValidateToken() fiber.Handler {
	//TODO implement me
	panic("implement me")
}

func (m *JWTMock) GenerateToken(user *entities.User) (*string, error) {
	arguments := m.Mock.Called(user)
	if arguments.Get(0) == nil {
		return nil, errors.New("user is nil")
	}

	result := fmt.Sprintf(`%s`, arguments.Get(0))
	return &result, nil
}
