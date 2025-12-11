package dto

type CSRRequest struct {
	CSR   []byte `json:"csr"`
	Email string `json:"email"`
}

type CertificateResponse struct {
	CsrId       int    `json:"csr_id"`
	Certificate []byte `json:"certificate"`
}
