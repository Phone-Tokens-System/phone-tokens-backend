package dto

type CSRRequest struct {
	CSR       string `json:"csr"`
	ServiceID string `json:"service_id"`
}
