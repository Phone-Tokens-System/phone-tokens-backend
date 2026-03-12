package model

import "time"

// тарифы смс и звонков
type Tariff struct {
	Service   ServiceType `json:"service"`
	Provider  string      `json:"provider"`
	Price     float64     `json:"price"`
	UpdatedAt time.Time   `json:"updated_at"`
}
