package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"phone_token_system/internal/model"
	"phone_token_system/internal/service/tokens"
	"phone_token_system/internal/service/users"
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
var _ tokens.Repository = (*Storage)(nil)

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

func (r *Storage) CreateToken(ctx context.Context, token *model.UserToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *Storage) GetTokenByID(ctx context.Context, id string) (*model.UserToken, error) {
	var token model.UserToken

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tokens.ErrNotFound
		}
		return nil, err
	}

	return &token, nil
}

func (r *Storage) UpdateToken(ctx context.Context, token *model.UserToken) error {
	result := r.db.WithContext(ctx).Model(&model.UserToken{}).Where("id = ?", token.ID).Update("expires_at", token.ExpiresAt)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return tokens.ErrNotFound
	}
	return nil
}

func (r *Storage) DeleteToken(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return tokens.ErrNotFound
	}
	return nil
}
