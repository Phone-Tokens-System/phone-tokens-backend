package tokens

import (
	"context"

	"phone-tokens/internal/model"

	"github.com/google/uuid"
)

// Repository defines storage operations for user tokens.
type Repository interface {
	CreateToken(ctx context.Context, token *model.UserToken) error
	GetTokenByID(ctx context.Context, id string) (*model.UserToken, error)
	UpdateToken(ctx context.Context, token *model.UserToken) error
	DeleteToken(ctx context.Context, id string) error
	GetUserIdFromToken(ctx context.Context, token string) (string, error)
	GetNumberFromUserId(ctx context.Context, userId string) (string, error)
	GetTokenByToken(ctx context.Context, token string) (*model.UserToken, error)
}

type IssueInput struct {
	UserID      string
	Name        string
	TTLSeconds  int64
	Permissions []model.TokenPermission
}

// Service defines token issuance use cases.
type Service interface {
	Issue(ctx context.Context, input IssueInput) (*model.UserToken, error)
	UpdateTTL(ctx context.Context, userID, tokenID string, ttlSeconds int64) (*model.UserToken, error)
	SetStatus(ctx context.Context, userID, tokenID string, status model.TokenStatus) (*model.UserToken, error)
	Delete(ctx context.Context, userID, tokenID string) error
	GetUserNumberFromToken(ctx context.Context, token string) (string, error)
	CheckTokenPermission(ctx context.Context, token string, agentId uuid.UUID, perm model.TokenPermission) (bool, error)
}
