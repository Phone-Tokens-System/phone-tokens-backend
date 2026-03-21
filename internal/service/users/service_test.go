package users

import (
	"context"
	"errors"
	"testing"

	"gorm.io/gorm"
	"phone-tokens/internal/model"
)

func TestRegisterDefaultsToUser(t *testing.T) {
	repo := newStubRepo()
	svc := NewService(repo, Config{JWTSecret: "secret", JWTExpiresInSec: 3600})

	user, err := svc.Register(context.Background(), RegisterRequest{
		Phone:    "100",
		Password: "password",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user.Role != model.RoleUser {
		t.Fatalf("expected role %q, got %q", model.RoleUser, user.Role)
	}
}

func TestRegisterAllowsAgent(t *testing.T) {
	repo := newStubRepo()
	svc := NewService(repo, Config{JWTSecret: "secret", JWTExpiresInSec: 3600})

	user, err := svc.Register(context.Background(), RegisterRequest{
		Phone:       "200",
		Password:    "password",
		Role:        model.RoleAgent,
		ServiceName: "svc",
		Email:       "agent@example.com",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user.Role != model.RoleAgent {
		t.Fatalf("expected role %q, got %q", model.RoleAgent, user.Role)
	}
	if repo.savedAgent == nil || repo.savedAgent.UserID != user.ID {
		t.Fatalf("expected agent to be created for user")
	}
}

func TestRegisterRejectsAdmin(t *testing.T) {
	repo := newStubRepo()
	svc := NewService(repo, Config{JWTSecret: "secret", JWTExpiresInSec: 3600})

	_, err := svc.Register(context.Background(), RegisterRequest{
		Phone:    "300",
		Password: "password",
		Role:     model.RoleAdmin,
	})
	if !errors.Is(err, ErrRoleNotAllowed) {
		t.Fatalf("expected ErrRoleNotAllowed, got %v", err)
	}
}

func TestRegisterAgentRequiresDetails(t *testing.T) {
	repo := newStubRepo()
	svc := NewService(repo, Config{JWTSecret: "secret", JWTExpiresInSec: 3600})

	if _, err := svc.Register(context.Background(), RegisterRequest{
		Phone:    "400",
		Password: "password",
		Role:     model.RoleAgent,
	}); !errors.Is(err, ErrAgentDetailsNeeded) {
		t.Fatalf("expected ErrAgentDetailsNeeded, got %v", err)
	}
}

type stubRepo struct {
	byPhone    map[string]*model.User
	savedAgent *model.Agent
}

func newStubRepo() *stubRepo {
	return &stubRepo{byPhone: make(map[string]*model.User)}
}

func (r *stubRepo) Save(ctx context.Context, entity interface{}) error {
	switch v := entity.(type) {
	case *model.User:
		r.byPhone[v.Phone] = v
		return nil
	default:
		return errors.New("unexpected entity type")
	}
}

func (r *stubRepo) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	user, ok := r.byPhone[phone]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *stubRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	for _, user := range r.byPhone {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, ErrNotFound
}

func (r *stubRepo) SaveAgent(ctx context.Context, agent *model.Agent) error {
	r.savedAgent = agent
	return nil
}

func (r *stubRepo) GetAgentByID(ctx context.Context, id string) (*model.Agent, error) {
	if r.savedAgent != nil && r.savedAgent.ID == id {
		return r.savedAgent, nil
	}
	return nil, ErrNotFound
}

func (r *stubRepo) GetAgentByUserID(ctx context.Context, userID string) (*model.Agent, error) {
	if r.savedAgent != nil && r.savedAgent.UserID == userID {
		return r.savedAgent, nil
	}
	return nil, ErrNotFound
}

func (r *stubRepo) GetNumberFromUserId(ctx context.Context, userId string) (string, error) {
	user, err := r.GetUserByID(ctx, userId)
	if err != nil {
		return "", err
	}
	return user.Phone, nil
}

func (r *stubRepo) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (r *stubRepo) UpdateAgent(ctx context.Context, tx *gorm.DB, agent *model.Agent) (*model.Agent, error) {
	r.savedAgent = agent
	return agent, nil
}

func (r *stubRepo) GetAgentForUpdate(ctx context.Context, tx *gorm.DB, id string) (*model.Agent, error) {
	return r.GetAgentByID(ctx, id)
}
