package model

import (
	"time"
)

type Agent struct {
	ID                 string    `json:"id" gorm:"type:uuid;primaryKey"`
	UserID             string    `json:"user_id" gorm:"type:uuid;not null;unique"`
	ServiceName        string    `json:"service_name" gorm:"not null"`
	Email              string    `json:"email" gorm:"not null"`
	Certificate        []byte    `json:"certificate" gorm:"type:bytea;not null;default:''"`
	CertificateRequest []byte    `json:"certificate_request" gorm:"type:bytea;not null;default:''"`
	Balance            float64   `json:"balance" gorm:"not null;default:0"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Agent) TableName() string {
	return "agents"
}
