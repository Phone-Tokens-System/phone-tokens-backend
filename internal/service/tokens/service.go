package tokens

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"

	"phone_token_system/internal/model"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Issue(ctx context.Context, userID string, ttlSeconds int64) (*model.UserToken, error) {
	if ttlSeconds <= 0 {
		return nil, errors.New("ttlSeconds must be greater than zero")
	}
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	value, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	token := &model.UserToken{
		ID:        uuid.NewString(),
		UserID:    userID,
		Token:     value,
		ExpiresAt: now.Add(time.Duration(ttlSeconds) * time.Second),
		CreatedAt: now,
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

	if err := s.repo.UpdateToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
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

func generateRandomToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
