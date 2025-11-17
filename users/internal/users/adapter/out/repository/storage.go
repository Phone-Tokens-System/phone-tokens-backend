package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"users/internal/users/model"
	"users/internal/users/service/users"
)

type PostgresRepository struct {
	db *gorm.DB
}

type Storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

var _ users.Repository = (*Storage)(nil)

func (r *Storage) Save(ctx context.Context, entity interface{}) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Storage) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Storage) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
