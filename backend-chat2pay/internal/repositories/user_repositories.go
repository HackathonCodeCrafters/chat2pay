package repositories

import (
	"chat2pay/internal/consts"
	"chat2pay/internal/entities"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindOneByEmail(ctx context.Context, email string) (*entities.User, error)
	FindOneById(ctx context.Context, id string) (*entities.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) FindOneByEmail(ctx context.Context, email string) (*entities.User, error) {
	user := entities.User{}
	err := r.db.WithContext(ctx).Where("email", email).First(&user).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindOneById(ctx context.Context, id string) (*entities.User, error) {
	user := entities.User{}
	err := r.db.WithContext(ctx).Where("id", id).First(&user).Error

	if err != nil {
		if err.Error() == consts.SqlNoRow {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
