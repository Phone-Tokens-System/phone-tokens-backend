package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TokenPermission string

const (
	TokenPermissionSMS   TokenPermission = "sms"
	TokenPermissionCalls TokenPermission = "calls"
)

type TokenStatus string

const (
	TokenStatusActive TokenStatus = "active"
	TokenStatusFrozen TokenStatus = "frozen"
)

type TokenPermissions []TokenPermission

// Value converts permissions to a Postgres-compatible text[] literal.
func (p TokenPermissions) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "{}", nil
	}

	parts := make([]string, len(p))
	for i, perm := range p {
		parts[i] = `"` + strings.ReplaceAll(string(perm), `"`, `\"`) + `"`
	}
	return "{" + strings.Join(parts, ",") + "}", nil
}

// Scan parses a Postgres text[] value into TokenPermissions.
func (p *TokenPermissions) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return p.fromString(v)
	case []byte:
		return p.fromString(string(v))
	default:
		return fmt.Errorf("cannot scan %T into TokenPermissions", src)
	}
}

func (p *TokenPermissions) fromString(s string) error {
	s = strings.Trim(s, "{}")
	if s == "" {
		*p = TokenPermissions{}
		return nil
	}

	parts := strings.Split(s, ",")
	perms := make(TokenPermissions, 0, len(parts))
	for _, part := range parts {
		part = strings.Trim(part, `"`)
		perms = append(perms, TokenPermission(part))
	}
	*p = perms
	return nil
}

type UserToken struct {
	ID          string           `json:"id" gorm:"type:uuid;primaryKey"`
	UserID      string           `json:"user_id" gorm:"type:uuid;index;not null"`
	Token       string           `json:"token" gorm:"not null"`
	AgentId     uuid.UUID        `json:"agent_id" gorm:"not null"`
	Name        string           `json:"name" gorm:"not null"`
	Permissions TokenPermissions `json:"permissions" gorm:"type:text[];not null"`
	Status      TokenStatus      `json:"status" gorm:"type:text;not null"`
	ExpiresAt   time.Time        `json:"expires_at" gorm:"not null;index"`
	CreatedAt   time.Time        `json:"created_at" gorm:"autoCreateTime"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
