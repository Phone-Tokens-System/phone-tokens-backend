package users

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"phone-tokens/internal/model"
)

var (
	ErrPhoneAlreadyUsed   = errors.New("phone already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrRoleNotAllowed     = errors.New("role must be user or agent")
	ErrAgentDetailsNeeded = errors.New("service_name and email are required for agent role")
)

type Config struct {
	JWTSecret       string
	JWTExpiresInSec int64
}

type service struct {
	repo   Repository
	config Config
}

func (s *service) GetAgentByID(ctx context.Context, id string) (*model.Agent, error) {
	return s.repo.GetAgentByID(ctx, id)
}

func (s *service) GetAgentByUserID(ctx context.Context, userID string) (*model.Agent, error) {
	return s.repo.GetAgentByUserID(ctx, userID)
}

func NewService(repo Repository, cfg Config) Service {
	return &service{
		repo:   repo,
		config: cfg,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*model.User, error) {
	role, err := sanitizeSelfAssignedRole(req.Role)
	if err != nil {
		return nil, err
	}

	if role == model.RoleAgent && (strings.TrimSpace(req.ServiceName) == "" || strings.TrimSpace(req.Email) == "") {
		return nil, ErrAgentDetailsNeeded
	}

	existing, err := s.repo.GetUserByPhone(ctx, req.Phone)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrPhoneAlreadyUsed
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &model.User{
		ID:           uuid.NewString(),
		Phone:        req.Phone,
		PasswordHash: string(hashed),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err = s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	if user.Role == model.RoleAgent {
		agent := &model.Agent{
			ID:                 uuid.NewString(),
			UserID:             user.ID,
			ServiceName:        req.ServiceName,
			Email:              req.Email,
			Certificate:        []byte{},
			CertificateRequest: []byte{},
			Balance:            0,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := s.repo.SaveAgent(ctx, agent); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *service) Authenticate(ctx context.Context, phone, password string) (string, *model.User, error) {
	user, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
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

func (s *service) GetAgentByID(ctx context.Context, id string) (*model.Agent, error) {
	return s.repo.GetAgentByID(ctx, id)
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

func sanitizeSelfAssignedRole(role model.Role) (model.Role, error) {
	if role == "" {
		return model.RoleUser, nil
	}

	if role != model.RoleUser && role != model.RoleAgent {
		return "", ErrRoleNotAllowed
	}

	return role, nil
}
