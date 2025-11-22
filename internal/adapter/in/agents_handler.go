package in

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/service"
)

type Handler struct {
	CertificateService *service.CertificateService
}

func NewHandler(certService *service.CertificateService) *Handler {
	return &Handler{certService}
}

func (h *Handler) AcceptCSRRequest(w http.ResponseWriter, r *http.Request) {
	var csr []byte

	err := json.NewDecoder(r.Body).Decode(&csr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	certificate, err := h.CertificateService.AcceptCertificate(r.Context(), csr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(certificate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShowCSRRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.CertificateService.GetCertificateRequests(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(requests)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) ApproveCSRRequest(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	request, err := h.CertificateService.ApproveCertificate(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
