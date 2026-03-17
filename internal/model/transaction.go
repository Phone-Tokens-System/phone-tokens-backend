package model

import "time"

// пополнения баланса
type Transaction struct {
	ID              int             `gorm:"id" json:"id"`
	AgentID         string          `gorm:"agent_id" json:"agent_id"`
	Amount          float64         `gorm:"amount" json:"amount"`
	Type            TranscationType `gorm:"type" json:"type"`
	Service         ServiceType     `gorm:"service" json:"service"`
	StripeSessionID string          `gorm:"stripe_session_id" json:"stripe_session_id"`
	CreatedAt       time.Time       `gorm:"created_at" json:"created_at"`
}

type TranscationType string

const (
	Debit  TranscationType = "debit"
	Credit TranscationType = "credit"
)

func (Transaction) TableName() string {
	return "billing_transactions"
}
