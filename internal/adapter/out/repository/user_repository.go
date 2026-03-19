package repository

import (
	"context"
	"errors"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/users"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *UserRepository) GetAgentByID(ctx context.Context, id string) (*model.Agent, error) {
	var agent model.Agent

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&agent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &agent, nil
}

func (r *UserRepository) GetAgentByUserID(ctx context.Context, userID string) (*model.Agent, error) {
	var agent model.Agent

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&agent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, users.ErrNotFound
		}
		return nil, err
	}

	return &agent, nil
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

func (r *UserRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (r *UserRepository) GetAgentForUpdate(ctx context.Context, tx *gorm.DB, id string) (*model.Agent, error) {
	var agent model.Agent
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&agent, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (r *UserRepository) UpdateAgent(ctx context.Context, tx *gorm.DB, agent *model.Agent) (*model.Agent, error) {
	if err := tx.Save(agent).Error; err != nil {
		return nil, err
	}
	return agent, nil
}
