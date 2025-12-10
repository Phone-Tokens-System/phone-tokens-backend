package users

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"phone-tokens/internal/model"
)

var (
	ErrPhoneAlreadyUsed   = errors.New("phone already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Config struct {
	JWTSecret       string
	JWTExpiresInSec int64
}

type service struct {
	repo   Repository
	config Config
}

func NewService(repo Repository, cfg Config) Service {
	return &service{
		repo:   repo,
		config: cfg,
	}
}

func (s *service) Register(ctx context.Context, phone, password string, role model.Role) (*model.User, error) {
	existing, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPhoneAlreadyUsed
	}

	if role == "" {
		role = model.RoleUser
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:           uuid.NewString(),
		Phone:        phone,
		PasswordHash: string(hashed),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err = s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) Authenticate(ctx context.Context, phone, password string) (string, *model.User, error) {
	user, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"phone": user.Phone,
		"role":  string(user.Role),
		"exp":   time.Now().Add(time.Duration(s.config.JWTExpiresInSec) * time.Second).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.config.JWTSecret))
}

var ErrNotFound = errors.New("user not found")
