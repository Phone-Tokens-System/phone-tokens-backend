package model

import (
	"time"
)

type Agent struct {
	ID                 string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	UserID             string    `json:"user_id" gorm:"column:user_id;type:uuid;not null;unique"`
	ServiceName        string    `json:"service_name" gorm:"column:service_name;not null"`
	Email              string    `json:"email" gorm:"column:email;not null"`
	Certificate        []byte    `json:"certificate" gorm:"column:certificate;type:bytea;not null;default:''"`
	CertificateRequest []byte    `json:"certificate_request" gorm:"column:certificate_request;type:bytea;not null;default:''"`
	Balance            float64   `json:"balance" gorm:"column:balance;not null;default:0"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Agent) TableName() string {
	return "agents"
}
