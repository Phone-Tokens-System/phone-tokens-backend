package repository

import (
	"context"
	"errors"
	"fmt"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/tokens"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) CreateToken(ctx context.Context, token *model.UserToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *TokenRepository) GetTokenByID(ctx context.Context, id string) (*model.UserToken, error) {
	var token model.UserToken

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tokens.ErrNotFound
		}
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) GetTokensByUserId(ctx context.Context, userId string) ([]model.UserToken, error) {
	var token []model.UserToken
	fmt.Println(userId)
	if err := r.db.WithContext(ctx).Find(&token, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	fmt.Println(token)
	return token, nil
}

func (r *TokenRepository) GetTokensByUserIdAndAgentId(ctx context.Context, userId, agentId string) ([]model.UserToken, error) {
	var token []model.UserToken
	fmt.Println(userId)
	if err := r.db.WithContext(ctx).Where("user_id = ? AND agent_id = ?", userId, agentId).Find(&token).Error; err != nil {
		return nil, err
	}
	fmt.Println(token)
	return token, nil
}

func (r *TokenRepository) GetTokensByAgentId(ctx context.Context, agentId string) ([]model.UserToken, error) {
	var result []model.UserToken
	if err := r.db.WithContext(ctx).Where("agent_id = ?", agentId).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *TokenRepository) UpdateToken(ctx context.Context, token *model.UserToken) (*model.UserToken, error) {
	var updatedToken model.UserToken
	result := r.db.WithContext(ctx).
		Model(&model.UserToken{}).
		Where("id = ?", token.ID).
		Updates(map[string]interface{}{
			"name":        token.Name,
			"permissions": token.Permissions,
			"status":      token.Status,
			"expires_at":  token.ExpiresAt,
			"agent_id":    token.AgentId,
		}).Scan(&updatedToken)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, tokens.ErrNotFound
	}

	return &updatedToken, nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserToken{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return tokens.ErrNotFound
	}
	return nil
}

func (r *TokenRepository) GetUserIdFromToken(ctx context.Context, token string) (string, error) {
	var user model.UserToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", model.ErrNotFound
		}
	}
	return user.UserID, nil
}

// TODO: remove
func (r *TokenRepository) GetNumberFromUserId(ctx context.Context, userId string) (string, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", model.ErrNotFound
		}
	}
	return user.Phone, nil
}

func (r *TokenRepository) GetTokenByToken(ctx context.Context, token string) (*model.UserToken, error) {
	var tokenObj model.UserToken
	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&tokenObj).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrNotFound
		}
	}
	return &tokenObj, nil
}
