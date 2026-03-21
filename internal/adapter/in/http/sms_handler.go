package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/sms"
	"phone-tokens/internal/service/users"
)

type SmsHandler struct {
	smsService *sms.SmsService
	userSvc    users.Service
}

func NewSmsHandler(smsService *sms.SmsService, userSvc users.Service) *SmsHandler {
	return &SmsHandler{smsService: smsService, userSvc: userSvc}
}

// SendSMS godoc
// @Summary Send an SMS
// @Description Sends an SMS to the specified phone number
// @Tags SMS
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.SmsRequest true "SMS request payload"
// @Success 200 {object} model.SmsResponse "Sent SMS details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/send [post]
func (h *SmsHandler) sendSMS(w http.ResponseWriter, req *http.Request) {
	var request model.SmsRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.userSvc != nil {
		claims, ok := req.Context().Value(userContextKey).(*UserClaims)
		if !ok || claims.UserID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		agent, err := h.userSvc.GetAgentByUserID(req.Context(), claims.UserID)
		if err != nil {
			http.Error(w, "failed to resolve agent profile", http.StatusInternalServerError)
			return
		}
		if agent != nil {
			request.AgentID = agent.ID
		}
	}

	sentSms, err := h.smsService.SendSms(req.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(sentSms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CheckStatus godoc
// @Summary Check SMS status
// @Description Returns the status of an SMS by ID
// @Tags SMS
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param smsId body dto.SmsId true "SMS ID payload"
// @Success 200 {object} model.SmsStatus "SMS status details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/status [post]
func (h *SmsHandler) checkStatus(w http.ResponseWriter, req *http.Request) {
	var smsId dto.SmsId
	err := json.NewDecoder(req.Body).Decode(&smsId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	status, err := h.smsService.GetSmsStatusPing(smsId.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetSmsList godoc
// @Summary Get all SMS
// @Security BearerAuth
// @Description Returns the list of all sent SMS
// @Tags SMS
// @Accept json
// @Produce json
// @Success 200 {array} model.SmsResponse "List of sent SMS"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/logs [get]
func (h *SmsHandler) getSmsList(w http.ResponseWriter, req *http.Request) {
	smsList, err := h.smsService.GetSmsList(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(smsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getSmsListByAgentId godoc
// @Summary get sms sent by agent by id
// @Security BearerAuth
// @Description Returns list of sms
// @Tags SMS
// @Accept json
// @Produce json
// @Param agentId path string true "Agent ID"
// @Success 200 {array} model.SmsResponse "SMS status details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/agents/{agentId} [get]
func (h *SmsHandler) getSmsListByAgentId(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("agentId")
	smsList, err := h.smsService.GetSmsListByAgentId(req.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(smsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getSmsListByToken godoc
// @Summary get sms sent to given token
// @Security BearerAuth
// @Description Returns list of sms
// @Tags SMS
// @Accept json
// @Produce json
// @Param token path string true "user token"
// @Success 200 {array} model.SmsResponse "SMS status details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/users/{token} [get]
func (h *SmsHandler) getSmsListByToken(w http.ResponseWriter, req *http.Request) {
	token := req.PathValue("token")
	fmt.Println(token)
	smsList, err := h.smsService.GetSmsByToken(req.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(smsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getSmsListFromProvider godoc
// @Summary get sms sent by provider
// @Security BearerAuth
// @Description Returns list of sms
// @Tags SMS
// @Accept json
// @Produce json
// @Success 200 {array} model.SmsResponse "SMS status details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/all [get]
func (h *SmsHandler) getSmsListFromProvider(w http.ResponseWriter, req *http.Request) {
	smsList, err := h.smsService.GetSmsListFromProvider(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(smsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SendSMSWithFilters godoc
// @Summary Send an SMS with filters for users
// @Description Sends an SMS to users who apply to filters
// @Tags SMS
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SmsFilterRequest true "SMS request payload"
// @Success 200 {array} []model.SmsResponse "Sent SMS details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/send_filtered [post]
func (h *SmsHandler) SendSmsWithFilters(w http.ResponseWriter, req *http.Request) {
	var smsReq dto.SmsFilterRequest
	if err := json.NewDecoder(req.Body).Decode(&smsReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	smsResp, err := h.smsService.SendSmsWithFilters(req.Context(), smsReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJSON(w, http.StatusOK, smsResp)
}
