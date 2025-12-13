package http

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/sms_service"
)

type SmsHandler struct {
	smsService *sms_service.SmsService
}

func NewSmsHandler(smsService *sms_service.SmsService) *SmsHandler {
	return &SmsHandler{smsService: smsService}
}

// SendSMS godoc
// @Summary Send an SMS
// @Description Sends an SMS to the specified phone number
// @Tags SMS
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

	sms, err := h.smsService.SendSms(req.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(sms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CheckStatus godoc
// @Summary Check SMS status
// @Description Returns the status of an SMS by ID
// @Tags SMS
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

	status, err := h.smsService.GetSmsStatus(smsId.Id)
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
// @Description Returns the list of all sent SMS
// @Tags SMS
// @Accept json
// @Produce json
// @Success 200 {array} model.SmsResponse "List of sent SMS"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/sms/list [get]
func (h *SmsHandler) getSmsList(w http.ResponseWriter, req *http.Request) {
	smsList, err := h.smsService.GetSmsList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(smsList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
