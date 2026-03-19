package tokens

import (
	"context"
	"errors"
	"testing"
	"time"

	"phone-tokens/internal/adapter/dto"

	"github.com/google/uuid"

	"phone-tokens/internal/model"
)

func TestServiceUpdateTTLSuccess(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token := &model.UserToken{
		ID:          "token-1",
		UserID:      "user-1",
		Token:       "value",
		Name:        "main",
		Permissions: model.TokenPermissions(DefaultTokenPermissions),
		Status:      model.TokenStatusActive,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
	}
	repo.tokens[token.ID] = token

	ttl := int64(120)
	before := time.Now().UTC()

	updated, err := svc.UpdateTTL(context.Background(), token.UserID, token.ID, ttl)
	if err != nil {
		t.Fatalf("UpdateTTL returned error: %v", err)
	}
	after := time.Now().UTC()

	lower := before.Add(time.Duration(ttl) * time.Second)
	upper := after.Add(time.Duration(ttl) * time.Second)
	if updated.ExpiresAt.Before(lower) || updated.ExpiresAt.After(upper) {
		t.Fatalf("ExpiresAt %v not in expected range [%v, %v]", updated.ExpiresAt, lower, upper)
	}

	stored := repo.tokens[token.ID]
	if stored == nil || !stored.ExpiresAt.Equal(updated.ExpiresAt) {
		t.Fatalf("repository not updated with new expiration")
	}
}

func TestServiceUpdateTTLForbidden(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token := &model.UserToken{
		ID:          "token-2",
		UserID:      "owner",
		Token:       "value",
		Name:        "secondary",
		Permissions: model.TokenPermissions(DefaultTokenPermissions),
		Status:      model.TokenStatusActive,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
	}
	repo.tokens[token.ID] = token

	_, err := svc.UpdateTTL(context.Background(), "other-user", token.ID, 60)
	if err := assertError(t, err, ErrForbidden); err != nil {
		t.Fatal(err)
	}
	if !repo.tokens[token.ID].ExpiresAt.Equal(token.ExpiresAt) {
		t.Fatalf("ExpiresAt changed for forbidden update")
	}
}

func TestServiceUpdateTTLNotFound(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	_, err := svc.UpdateTTL(context.Background(), "user", "missing", 60)
	if err := assertError(t, err, ErrNotFound); err != nil {
		t.Fatal(err)
	}
}

func TestServiceUpdateTTLValidation(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	if _, err := svc.UpdateTTL(context.Background(), "", "id", 60); err == nil {
		t.Fatalf("expected error for missing userID")
	}
	if _, err := svc.UpdateTTL(context.Background(), "user", "", 60); err == nil {
		t.Fatalf("expected error for missing tokenID")
	}
	if _, err := svc.UpdateTTL(context.Background(), "user", "id", 0); err == nil {
		t.Fatalf("expected error for invalid ttl")
	}
}

func TestServiceDeleteSuccess(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token := &model.UserToken{
		ID:          "token-3",
		UserID:      "user-1",
		Token:       "value",
		Name:        "to-delete",
		Permissions: model.TokenPermissions(DefaultTokenPermissions),
		Status:      model.TokenStatusActive,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
	}
	repo.tokens[token.ID] = token

	if err := svc.Delete(context.Background(), token.UserID, token.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if _, ok := repo.tokens[token.ID]; ok {
		t.Fatalf("token was not removed from repository")
	}
}

func TestServiceDeleteForbidden(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token := &model.UserToken{
		ID:          "token-4",
		UserID:      "owner",
		Token:       "value",
		Name:        "locked",
		Permissions: model.TokenPermissions(DefaultTokenPermissions),
		Status:      model.TokenStatusActive,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
	}
	repo.tokens[token.ID] = token

	if err := assertError(t, svc.Delete(context.Background(), "other", token.ID), ErrForbidden); err != nil {
		t.Fatal(err)
	}
	if _, ok := repo.tokens[token.ID]; !ok {
		t.Fatalf("token should remain after forbidden delete")
	}
}

func TestServiceDeleteNotFound(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	if err := assertError(t, svc.Delete(context.Background(), "user", "missing"), ErrNotFound); err != nil {
		t.Fatal(err)
	}
}

func TestServiceDeleteValidation(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	if err := svc.Delete(context.Background(), "", "id"); err == nil {
		t.Fatalf("expected error for missing userID")
	}
	if err := svc.Delete(context.Background(), "user", ""); err == nil {
		t.Fatalf("expected error for missing tokenID")
	}
}

func TestServiceIssueDefaults(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	ttl := int64(120)
	before := time.Now().UTC()
	token, err := svc.Issue(context.Background(), IssueInput{
		UserID:     "user-issue",
		TTLSeconds: ttl,
		AgentId:    uuid.Nil.String(),
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}
	if token.Name != DefaultTokenName {
		t.Fatalf("expected default name %q, got %q", DefaultTokenName, token.Name)
	}
	if len(token.Permissions) != len(DefaultTokenPermissions) {
		t.Fatalf("expected default permissions, got %v", token.Permissions)
	}
	if token.Status != model.TokenStatusActive {
		t.Fatalf("expected status active, got %s", token.Status)
	}
	if token.ExpiresAt.Before(before.Add(time.Duration(ttl) * time.Second)) {
		t.Fatalf("unexpected expires_at: %v", token.ExpiresAt)
	}
}

func TestServiceIssueCustomPermissions(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token, err := svc.Issue(context.Background(), IssueInput{
		UserID:     "user-perms",
		TTLSeconds: 60,
		Name:       "calls-only",
		AgentId:    uuid.Nil.String(),
		Permissions: []model.TokenPermission{
			model.TokenPermissionCalls,
			model.TokenPermissionCalls,
		},
	})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	if len(token.Permissions) != 1 || token.Permissions[0] != model.TokenPermissionCalls {
		t.Fatalf("unexpected permissions: %v", token.Permissions)
	}
}

func TestServiceSetStatus(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token := &model.UserToken{
		ID:          "token-5",
		UserID:      "owner",
		Token:       "value",
		Name:        "freeze-me",
		Permissions: model.TokenPermissions(DefaultTokenPermissions),
		Status:      model.TokenStatusActive,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
	}
	repo.tokens[token.ID] = token

	updated, err := svc.SetStatus(context.Background(), token.UserID, token.ID, model.TokenStatusFrozen)
	if err != nil {
		t.Fatalf("SetStatus returned error: %v", err)
	}
	if updated.Status != model.TokenStatusFrozen {
		t.Fatalf("expected frozen status, got %s", updated.Status)
	}
	if repo.tokens[token.ID].Status != model.TokenStatusFrozen {
		t.Fatalf("repository not updated")
	}
}

func TestBindAgentRequiresOwnership(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token := &model.UserToken{
		ID:          "token-bind",
		UserID:      "owner",
		Token:       "value",
		Name:        "bind-me",
		Permissions: model.TokenPermissions(DefaultTokenPermissions),
		Status:      model.TokenStatusActive,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
	}
	repo.tokens[token.ID] = token

	_, err := svc.BingAgentToTokenByName(context.Background(), "other-user", dto.BindTokenRequest{
		AgentId:   uuid.NewString(),
		TokenName: token.Token,
	})
	if err := assertError(t, err, ErrForbidden); err != nil {
		t.Fatal(err)
	}
}

func TestBindAgentSuccess(t *testing.T) {
	repo := newMemoryRepo()
	svc := NewService(repo)

	token := &model.UserToken{
		ID:          "token-bind-success",
		UserID:      "owner",
		Token:       "value",
		Name:        "bind-me",
		Permissions: model.TokenPermissions(DefaultTokenPermissions),
		Status:      model.TokenStatusActive,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
	}
	repo.tokens[token.ID] = token

	agentID := uuid.NewString()

	updated, err := svc.BingAgentToTokenByName(context.Background(), token.UserID, dto.BindTokenRequest{
		AgentId:   agentID,
		TokenName: token.Token,
	})
	if err != nil {
		t.Fatalf("BingAgentToTokenByName returned error: %v", err)
	}
	if updated.AgentId.String() != agentID {
		t.Fatalf("expected agent id %s, got %s", agentID, updated.AgentId)
	}
	if repo.tokens[token.ID].AgentId.String() != agentID {
		t.Fatalf("repository not updated with agent id")
	}
}

type memoryRepo struct {
	tokens      map[string]*model.UserToken
	phonesByUID map[string]string
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{
		tokens:      make(map[string]*model.UserToken),
		phonesByUID: make(map[string]string),
	}
}

func (r *memoryRepo) CreateToken(_ context.Context, token *model.UserToken) error {
	r.tokens[token.ID] = token
	return nil
}

func (r *memoryRepo) GetTokenByID(_ context.Context, id string) (*model.UserToken, error) {
	token, ok := r.tokens[id]
	if !ok {
		return nil, ErrNotFound
	}
	return token, nil
}

func (r *memoryRepo) UpdateToken(_ context.Context, token *model.UserToken) (*model.UserToken, error) {
	if _, ok := r.tokens[token.ID]; !ok {
		return nil, ErrNotFound
	}
	r.tokens[token.ID] = token
	return token, nil
}

func (r *memoryRepo) DeleteToken(_ context.Context, id string) error {
	if _, ok := r.tokens[id]; !ok {
		return ErrNotFound
	}
	delete(r.tokens, id)
	return nil
}

func (r *memoryRepo) GetUserIdFromToken(_ context.Context, token string) (string, error) {
	for _, t := range r.tokens {
		if t.Token == token {
			return t.UserID, nil
		}
	}
	return "", ErrNotFound
}

func (r *memoryRepo) GetNumberFromUserId(_ context.Context, userId string) (string, error) {
	if phone, ok := r.phonesByUID[userId]; ok {
		return phone, nil
	}
	return "", ErrNotFound
}

func (r *memoryRepo) GetTokenByToken(_ context.Context, token string) (*model.UserToken, error) {
	for _, t := range r.tokens {
		if t.Token == token {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

func (r *memoryRepo) GetTokensByUserId(_ context.Context, userId string) ([]model.UserToken, error) {
	result := make([]model.UserToken, 0)
	for _, t := range r.tokens {
		if t.UserID == userId {
			result = append(result, *t)
		}
	}
	return result, nil
}

func assertError(t *testing.T, err error, expected error) error {
	t.Helper()
	if err == nil {
		return errors.New("expected error, got nil")
	}
	if !errors.Is(err, expected) {
		return errors.New("unexpected error: " + err.Error())
	}
	return nil
}
