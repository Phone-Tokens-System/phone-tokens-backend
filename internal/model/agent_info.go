package model

type ExternalAgentInfo struct {
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
	CertificatePem string `json:"certificate_pem"`
	IsActive       bool   `json:"is_active"`
}
