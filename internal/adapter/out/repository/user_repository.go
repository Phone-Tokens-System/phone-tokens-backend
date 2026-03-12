package repository

import (
	"context"
	"errors"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/users"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *UserRepository) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) SaveAgent(ctx context.Context, agent *model.Agent) error {
	return r.db.WithContext(ctx).Save(agent).Error
}

func (r *UserRepository) GetNumberFromUserId(ctx context.Context, userId string) (string, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", users.ErrNotFound
		}
	}
	return user.Phone, nil
}
