package model

import (
	"time"
)

type UserProfile struct {
	UserId    string    `json:"user_id" gorm:"primaryKey;type:uuid"`
	BirthDate time.Time `json:"birth_date" gorm:"type:VARCHAR(255)"`
	Age       int       `json:"age" gorm:"type:INTEGER"`
	Gender    string    `json:"gender" gorm:"type:VARCHAR(10)"`
	Country   string    `json:"country" gorm:"type:VARCHAR(255)"`
	Region    string    `json:"region" gorm:"type:VARCHAR(255)"`
	City      string    `json:"city" gorm:"type:VARCHAR(255)"`
	Education string    `json:"education" gorm:"type:VARCHAR(255)"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
