package in

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/service/certificates"
	"strconv"
)

type AgentHandler struct {
	CertificateService *certificates.CertificateService
}

func NewHandler(certService *certificates.CertificateService) *AgentHandler {
	return &AgentHandler{certService}
}

// AcceptCSRRequest godoc
// @Summary Accept a new CSR request
// @Description Accepts a Certificate Signing Request from an agent
// @Tags CSR
// @Accept json
// @Produce json
// @Param csr body dto.CSRRequest true "CSR Request payload"
// @Success 201 {object} dto.CertificateResponse "Signed certificate"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /csr [post]
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

// ShowCSRRequests godoc
// @Summary List all pending CSR requests
// @Description Returns all pending Certificate Signing Requests
// @Tags CSR
// @Accept json
// @Produce json
// @Success 200 {array} dto.CSRRequest "List of CSR requests"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/csr [get]
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

// ApproveCSRRequest godoc
// @Summary Approve a CSR request
// @Description Approves a CSR request by its ID and signs the certificate
// @Tags CSR
// @Accept json
// @Produce json
// @Param id query int true "CSR Request ID"
// @Success 200 {object} model.CsrRequest "Signed certificate"
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /admin/csr/approve [post]
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

	err = json.NewEncoder(w).Encode(*request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetSignedCertificate godoc
// @Summary Get signed certificate by CSR ID
// @Description Returns the signed certificate for the given CSR ID
// @Tags CSR
// @Accept json
// @Produce json
// @Param id query int true "CSR ID"
// @Success 200 {object} dto.CertificateResponse "Signed certificate"
// @Failure 400 {object} map[string]string "Invalid CSR ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /csr/signed [get]
func (h *AgentHandler) GetSignedCertificate(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cert, err := h.CertificateService.GetSignedCertificateByCsrID(r.Context(), idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := dto.CertificateResponse{
		Certificate: cert,
		CsrId:       idInt,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
