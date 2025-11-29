package dto

type CSRRequest struct {
	CSR   []byte `json:"csr"`
	Email string `json:"email"`
}
