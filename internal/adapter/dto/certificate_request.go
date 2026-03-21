package dto

type CSRRequest struct {
	AgentID string `json:"agent_id;omitempty"`
	CSR     string `json:"csr"`
	Email   string `json:"email"`
}

type CertificateResponse struct {
	CsrId       int    `json:"csr_id"`
	Certificate string `json:"certificate"`
}
