package http

import (
	"encoding/json"
	"io"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/service/certificates"
	"phone-tokens/internal/service/sms"
	"phone-tokens/internal/service/users"
	"strconv"
)

type AgentHandler struct {
	CertificateService *certificates.CertificateService
	UserService        users.Service
	SmsService         *sms.SmsService
}

func NewAgentHandler(certService *certificates.CertificateService, userService users.Service, smsService *sms.SmsService) *AgentHandler {
	return &AgentHandler{certService, userService, smsService}
}

// UploadCSR godoc
// @Summary Upload a Certificate Signing Request (CSR) file
// @Description Accepts a CSR from an agent as a file upload
// @Tags CSR
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param csr formData file true "CSR file"
// @Param email formData string true "Email of requester"
// @Success 201 {object} dto.CertificateResponse "Signed certificate"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/csr/upload [post]
func (h *AgentHandler) UploadCSR(w http.ResponseWriter, r *http.Request) {
	agent, err := GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	agentReal, err := h.UserService.GetAgentByUserID(r.Context(), agent.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	agentID := agentReal.ID

	err = r.ParseMultipartForm(1 << 20)
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

	csr := dto.CSRRequest{
		CSR:     string(csrBytes),
		Email:   email,
		AgentID: agentID,
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
// @Security BearerAuth
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

// SeeSMSLogs godoc
// @Summary Get sms logs for given agent
// @Description Returns all sms logs with given agent id
// @Tags Agents
// @Accept json
// @Produce json
// @Param id query string true "agent ID"
// @Success 200 {object} dto.SmsLog "Signed certificate"
// @Failure 400 {object} map[string]string "Invalid agent ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/agents/logs [get]
func (h *AgentHandler) SeeSMSLogs(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	smsResp, err := h.SmsService.GetSmsListByAgentId(r.Context(), id)
	if err != nil {
		return
	}
	resp := make([]dto.SmsLog, len(smsResp))
	for _, sm := range smsResp {
		resp = append(resp, *dto.ToSmsLog(sm))
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
