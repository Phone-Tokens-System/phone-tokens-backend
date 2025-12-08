package in

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/sms_service/model"
	"phone-tokens/internal/sms_service/service"
)

type smsHandler struct {
	smsService *service.SmsService
}

func (h *smsHandler) sendSMS(w http.ResponseWriter, req *http.Request) {
	var request model.SmsRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sms, err := h.smsService.SendSms(request)
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

func (h *smsHandler) checkStatus(w http.ResponseWriter, req *http.Request) {
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
