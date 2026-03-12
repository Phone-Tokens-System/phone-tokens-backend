package model

import "time"

// Usage запись истории операций смс/звонков
type Usage struct {
	ID          int         `gorm:"id" json:"id"`
	AgentID     string      `gorm:"agent_id" json:"agent_id"`
	PhoneNumber string      `gorm:"phone_number" json:"phone_number"`
	Service     ServiceType `gorm:"service_type" json:"service_type"`
	Units       int         `gorm:"units" json:"units"`
	Cost        float64     `gorm:"cost" json:"cost"`
	CreatedAt   time.Time   `gorm:"created_at" json:"created_at"`
}

type ServiceType string

const (
	SMS  ServiceType = "SMS"
	Call ServiceType = "Call"
)

func (Usage) TableName() string {
	return "usage"
}
