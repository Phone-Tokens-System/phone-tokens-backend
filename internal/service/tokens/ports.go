package tokens

import (
	"context"

	"phone_token_system/internal/model"
)

// Repository defines storage operations for user tokens.
type Repository interface {
	CreateToken(ctx context.Context, token *model.UserToken) error
}

// Service defines token issuance use cases.
type Service interface {
	Issue(ctx context.Context, userID string, ttlSeconds int64) (*model.UserToken, error)
}

