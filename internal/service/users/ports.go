package users

import (
	"context"

	"phone-tokens/internal/model"
)

type Repository interface {
	Save(ctx context.Context, entity interface{}) error
	GetUserByPhone(ctx context.Context, phone string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	SaveAgent(ctx context.Context, agent *model.Agent) error
}

type RegisterRequest struct {
	Phone       string     `json:"phone"`
	Password    string     `json:"password"`
	Role        model.Role `json:"role"`
	ServiceName string     `json:"service_name"`
	Email       string     `json:"email"`
}

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*model.User, error)
	Authenticate(ctx context.Context, phone, password string) (string, *model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}
