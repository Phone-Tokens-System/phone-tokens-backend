package in

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/certificates/service"
	"strconv"
)

type AgentHandler struct {
	CertificateService *service.CertificateService
}

func NewHandler(certService *service.CertificateService) *AgentHandler {
	return &AgentHandler{certService}
}

func (h *AgentHandler) AcceptCSRRequest(w http.ResponseWriter, r *http.Request) {
	var csr dto.CSRRequest

	err := json.NewDecoder(r.Body).Decode(&csr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	certificate, err := h.CertificateService.AcceptCertificateRequest(r.Context(), csr)
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

func (h *AgentHandler) ShowCSRRequests(w http.ResponseWriter, r *http.Request) {
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

func (h *AgentHandler) ApproveCSRRequest(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idInt, err := strconv.Atoi(id)
	request, err := h.CertificateService.ApproveCertificateRequest(r.Context(), idInt)
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
