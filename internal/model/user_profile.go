package model

import "github.com/google/uuid"

type UserProfile struct {
	UserId    uuid.UUID `json:"user_id" gorm:"primaryKey;type:uuid"`
	BirthDate string    `json:"birth_date" gorm:"type:VARCHAR(255)"`
	Gender    string    `json:"gender" gorm:"type:VARCHAR(10)"`
	Country   string    `json:"country" gorm:"type:VARCHAR(255)"`
	Region    string    `json:"region" gorm:"type:VARCHAR(255)"`
	City      string    `json:"city" gorm:"type:VARCHAR(255)"`
	Education string    `json:"education" gorm:"type:VARCHAR(255)"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
