package users

import (
	"context"
	"errors"
	"testing"

	"phone-tokens/internal/model"
)

func TestRegisterDefaultsToUser(t *testing.T) {
	repo := newStubRepo()
	svc := NewService(repo, Config{JWTSecret: "secret", JWTExpiresInSec: 3600})

	user, err := svc.Register(context.Background(), "100", "password", "")
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

	user, err := svc.Register(context.Background(), "200", "password", model.RoleAgent)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user.Role != model.RoleAgent {
		t.Fatalf("expected role %q, got %q", model.RoleAgent, user.Role)
	}
}

func TestRegisterRejectsAdmin(t *testing.T) {
	repo := newStubRepo()
	svc := NewService(repo, Config{JWTSecret: "secret", JWTExpiresInSec: 3600})

	_, err := svc.Register(context.Background(), "300", "password", model.RoleAdmin)
	if !errors.Is(err, ErrRoleNotAllowed) {
		t.Fatalf("expected ErrRoleNotAllowed, got %v", err)
	}
}

type stubRepo struct {
	byPhone map[string]*model.User
}

func newStubRepo() *stubRepo {
	return &stubRepo{byPhone: make(map[string]*model.User)}
}

func (r *stubRepo) Save(ctx context.Context, entity interface{}) error {
	user, ok := entity.(*model.User)
	if !ok {
		return errors.New("unexpected entity type")
	}
	r.byPhone[user.Phone] = user
	return nil
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
