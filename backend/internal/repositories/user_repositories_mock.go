package repositories

import (
	"chat2pay/internal/entities"
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (u *UserRepositoryMock) FindOneByEmail(ctx context.Context, email string) (*entities.User, error) {
	arguments := u.Mock.Called(email)
	if arguments.Get(0) == nil {
		return nil, nil
	}

	if arguments.Get(0) == "error banh" {
		return nil, errors.New("error")
	}

	user := arguments.Get(0).(entities.User)
	return &user, nil
}

func (u *UserRepositoryMock) FindOneById(ctx context.Context, id string) (*entities.User, error) {
	arguments := u.Mock.Called(id)
	if arguments.Get(0) == nil {
		return nil, nil
	}

	if arguments.Get(0) == "error banh" {
		return nil, errors.New("error")
	}

	user := arguments.Get(0).(entities.User)
	return &user, nil
}
