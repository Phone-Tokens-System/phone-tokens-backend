package tokens

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"phone-tokens/internal/adapter/dto"
	"strings"
	"time"

	"github.com/google/uuid"

	"phone-tokens/internal/model"
)

type service struct {
	repo Repository
}

const DefaultTokenName = "default token"

var DefaultTokenPermissions = []model.TokenPermission{
	model.TokenPermissionSMS,
	model.TokenPermissionCalls,
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Issue(ctx context.Context, input IssueInput) (*model.UserToken, error) {
	if input.TTLSeconds <= 0 {
		return nil, errors.New("ttlSeconds must be greater than zero")
	}
	if input.UserID == "" {
		return nil, errors.New("userID is required")
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		name = DefaultTokenName
	}

	perms, err := normalizePermissions(input.Permissions)
	if err != nil {
		return nil, err
	}

	value, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	token := &model.UserToken{
		ID:          uuid.NewString(),
		UserID:      input.UserID,
		Token:       value,
		Name:        name,
		Permissions: model.TokenPermissions(perms),
		Status:      model.TokenStatusActive,
		ExpiresAt:   now.Add(time.Duration(input.TTLSeconds) * time.Second),
		CreatedAt:   now,
	}

	if err := s.repo.CreateToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *service) UpdateTTL(ctx context.Context, userID, tokenID string, ttlSeconds int64) (*model.UserToken, error) {
	if ttlSeconds <= 0 {
		return nil, errors.New("ttlSeconds must be greater than zero")
	}
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	if tokenID == "" {
		return nil, errors.New("tokenID is required")
	}

	token, err := s.repo.GetTokenByID(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	if token.UserID != userID {
		return nil, ErrForbidden
	}

	now := time.Now().UTC()
	token.ExpiresAt = now.Add(time.Duration(ttlSeconds) * time.Second)

	updatedToken, err := s.repo.UpdateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return updatedToken, nil
}

func (s *service) SetStatus(ctx context.Context, userID, tokenID string, status model.TokenStatus) (*model.UserToken, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	if tokenID == "" {
		return nil, errors.New("tokenID is required")
	}
	if err := validateStatus(status); err != nil {
		return nil, err
	}

	token, err := s.repo.GetTokenByID(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	if token.UserID != userID {
		return nil, ErrForbidden
	}

	token.Status = status

	updatedToken, err := s.repo.UpdateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return updatedToken, nil
}

func (s *service) Delete(ctx context.Context, userID, tokenID string) error {
	if userID == "" {
		return errors.New("userID is required")
	}
	if tokenID == "" {
		return errors.New("tokenID is required")
	}

	token, err := s.repo.GetTokenByID(ctx, tokenID)
	if err != nil {
		return err
	}
	if token.UserID != userID {
		return ErrForbidden
	}

	return s.repo.DeleteToken(ctx, tokenID)
}

func normalizePermissions(perms []model.TokenPermission) ([]model.TokenPermission, error) {
	if len(perms) == 0 {
		return append([]model.TokenPermission(nil), DefaultTokenPermissions...), nil
	}

	seen := make(map[model.TokenPermission]struct{})
	result := make([]model.TokenPermission, 0, len(perms))

	for _, perm := range perms {
		if !isValidPermission(perm) {
			return nil, fmt.Errorf("%w: %s", ErrInvalidPermission, perm)
		}
		if _, ok := seen[perm]; ok {
			continue
		}
		seen[perm] = struct{}{}
		result = append(result, perm)
	}

	return result, nil
}

func isValidPermission(perm model.TokenPermission) bool {
	switch perm {
	case model.TokenPermissionSMS, model.TokenPermissionCalls:
		return true
	default:
		return false
	}
}

func validateStatus(status model.TokenStatus) error {
	switch status {
	case model.TokenStatusActive, model.TokenStatusFrozen:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}
}

func generateRandomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *service) GetUserNumberFromToken(ctx context.Context, token string) (string, error) {
	userId, err := s.repo.GetUserIdFromToken(ctx, token)
	if err != nil {
		return "", err
	}
	number, err := s.repo.GetNumberFromUserId(ctx, userId)
	if err != nil {
		return "", err
	}
	return number, nil
}

func (s *service) CheckTokenPermission(ctx context.Context, token string, agentId uuid.UUID, perm model.TokenPermission) (bool, error) {
	tokenObj, err := s.repo.GetTokenByToken(ctx, token)
	if err != nil {
		return false, err
	}

	if tokenObj.AgentId != agentId {
		return false, ErrForbidden
	}

	for _, tokenPerm := range tokenObj.Permissions {
		if tokenPerm == perm {
			return true, nil
		}
	}
	return false, nil
}

func (s *service) BingAgentToTokenByName(ctx context.Context, userID string, request dto.BindTokenRequest) (*model.UserToken, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}
	tokenObj, err := s.repo.GetTokenByToken(ctx, request.TokenName)
	if err != nil {
		return nil, err
	}
	if tokenObj.UserID != userID {
		return nil, ErrForbidden
	}
	uid, err := uuid.Parse(request.AgentId)
	if err != nil {
		return nil, err
	}
	tokenObj.AgentId = uid
	token, err := s.repo.UpdateToken(ctx, tokenObj)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *service) GetTokensByUser(ctx context.Context, userID string) ([]model.UserToken, error) {
	tokens, err := s.repo.GetTokensByUserId(ctx, userID)
	fmt.Println(tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
