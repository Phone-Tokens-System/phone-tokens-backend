package model

import "time"

// Usage запись истории операций смс/звонков
type Usage struct {
	ID          int         `json:"id"`
	AgentID     string      `json:"agent_id"`
	PhoneNumber string      `json:"phone_number"`
	Service     ServiceType `json:"service_type"`
	Units       int         `json:"units"`
	Cost        float64     `json:"cost"`
	CreatedAt   time.Time   `json:"created_at"`
}

type ServiceType string

const (
	SMS  ServiceType = "SMS"
	Call ServiceType = "Call"
)
