package model

import "github.com/google/uuid"

type ExternalAgentInfo struct {
	CsrID          int       `json:"csr_id"`
	ID             uuid.UUID `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ServiceName    string    `json:"service_name"`
	OrganizationID string    `json:"organization_id" gorm:"size:500; UNIQUE_INDEX;NOT NULL;"`
	Email          string    `json:"email" gorm:"size:255;NOT NULL;"`
	CertificatePem []byte    `json:"certificate_pem" gorm:"size:65535;NOT NULL;"`
	IsActive       bool      `json:"is_active" gorm:"size:1;NOT NULL;"`
}

func (ExternalAgentInfo) TableName() string {
	return "agent_info"
}
