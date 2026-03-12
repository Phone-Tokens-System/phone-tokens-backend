package http

import (
	"encoding/json"
	"io"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/service/certificates"
	"strconv"
)

type AgentHandler struct {
	CertificateService *certificates.CertificateService
}

func NewAgentHandler(certService *certificates.CertificateService) *AgentHandler {
	return &AgentHandler{certService}
}

// UploadCSR godoc
// @Summary Upload a Certificate Signing Request (CSR) file
// @Description Accepts a CSR from an agent as a file upload
// @Tags CSR
// @Accept multipart/form-data
// @Produce json
// @Param csr formData file true "CSR file"
// @Param email formData string true "Email of requester"
// @Success 201 {object} dto.CertificateResponse "Signed certificate"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/csr/upload [post]
func (h *AgentHandler) UploadCSR(w http.ResponseWriter, r *http.Request) {
	// Ограничим размер файла (например, 1MB)
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		return
	}

	file, _, err := r.FormFile("csr")
	if err != nil {
		http.Error(w, "CSR file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	csrBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read CSR file", http.StatusInternalServerError)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	// Передаем bytes в сервис
	csr := dto.CSRRequest{
		CSR:   string(csrBytes),
		Email: email,
	}

	certificate, err := h.CertificateService.AcceptCertificateRequest(r.Context(), csr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, certificate)
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
// @Router /api/v1/csr/signed [get]
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
		Certificate: string(cert),
		CsrId:       idInt,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
