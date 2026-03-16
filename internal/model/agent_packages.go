package model

import (
	"time"
)

type AgentPackages struct {
	Id         int         `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	AgentId    string      `json:"agent_id" gorm:"type:uuid;not null"`
	PackageId  string      `json:"package_id" gorm:"type:uuid;not null"`
	Service    ServiceType `json:"service_type" gorm:"type:varchar(50);not null"`
	Status     string      `json:"status" gorm:"type:varchar(50); not null"`
	UnitsTotal int64       `json:"units_total" gorm:"type:int"`
	UnitsUsed  int64       `json:"units_used" gorm:"type:int"`
	ExpiresAt  time.Time   `json:"expires_at" gorm:"type:datetime"`
}

func (AgentPackages) TableName() string {
	return "agent_packages"
}
