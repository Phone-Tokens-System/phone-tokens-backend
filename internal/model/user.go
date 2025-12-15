package model

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleAgent Role = "agent"
)

type User struct {
	ID           string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Phone        string    `json:"phone" gorm:"column:phone;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;not null"`
	Role         Role      `json:"role" gorm:"column:role;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
