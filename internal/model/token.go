package model

import "time"

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

type UserToken struct {
	ID          string            `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	UserID      string            `json:"user_id" gorm:"column:user_id;type:uuid;index;not null"`
	Token       string            `json:"token" gorm:"column:token;not null"`
	Name        string            `json:"name" gorm:"column:name;not null"`
	Permissions []TokenPermission `json:"permissions" gorm:"column:permissions;type:text[];not null"`
	Status      TokenStatus       `json:"status" gorm:"column:status;type:text;not null"`
	ExpiresAt   time.Time         `json:"expires_at" gorm:"column:expires_at;not null;index"`
	CreatedAt   time.Time         `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
