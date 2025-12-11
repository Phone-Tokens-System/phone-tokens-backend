package model

import "github.com/google/uuid"

type ExternalAgentInfo struct {
	CsrID          int       `json:"csr_id" gorm:"column:csr_id;primary_key"`
	ID             uuid.UUID `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ServiceName    string    `json:"service_name" gorm:"column:service_name"`
	OrganizationID string    `json:"organization_id" gorm:"column:organization_id; size:500; UNIQUE_INDEX;NOT NULL;"`
	Email          string    `json:"email" gorm:"column:email; size:255;NOT NULL;"`
	CertificatePem []byte    `json:"certificate_pem" gorm:"column:certificate_pem; size:65535;NOT NULL;"`
	IsActive       bool      `json:"is_active" gorm:"column:is_active; size:1;NOT NULL;"`
}

func (ExternalAgentInfo) TableName() string {
	return "agent_info"
}
