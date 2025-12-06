package model

import "time"

type UserToken struct {
	ID        string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	UserID    string    `json:"user_id" gorm:"column:user_id;type:uuid;index;not null"`
	Token     string    `json:"token" gorm:"column:token;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"column:expires_at;not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}

