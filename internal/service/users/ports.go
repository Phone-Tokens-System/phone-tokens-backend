package users

import (
	"context"

	"phone-tokens/internal/model"
)

type Repository interface {
	Save(ctx context.Context, entity interface{}) error
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

type Service interface {
	Register(ctx context.Context, phone, password string, role model.Role) (*model.User, error)
	Authenticate(ctx context.Context, phone, password string) (string, *model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}
