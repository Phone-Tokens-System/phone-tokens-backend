package in

import "phone-tokens/internal/service"

type Handler struct {
	CertificateService *service.CertificateService
}

func NewHandler(certService *service.CertificateService) *Handler {
	return &Handler{certService}
}

func (h *Handler) AcceptRequest() {

}
