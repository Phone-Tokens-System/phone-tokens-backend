package model

import "time"

// пополнения баланса
type Transaction struct {
	ID        int             `json:"id"`
	AgentID   string          `json:"agent_id"`
	Amount    float64         `json:"amount"`
	Type      TranscationType `json:"type"`
	Service   ServiceType     `json:"service"`
	CreatedAt time.Time       `json:"created_at"`
}

type TranscationType string

const (
	Debit  TranscationType = "debit"
	Credit TranscationType = "credit"
)
