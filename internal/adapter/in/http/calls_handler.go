package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"phone-tokens/internal/service/calls"
	"phone-tokens/internal/service/calls/novofon"
)

type CallHandler struct {
	service calls.Service
}

func NewCallHandler(service calls.Service) *CallHandler {
	return &CallHandler{service: service}
}

type connectCallRequest struct {
	ClientNumber string `json:"client_number"`
	UserNumber   string `json:"user_number"`
	ProxyNumber  string `json:"proxy_number,omitempty"`
	Predicted    bool   `json:"predicted"`
}

// ConnectClientWithUserViaProxy godoc
// @Summary Connect client and user (optionally via proxy number)
// @Security BearerAuth
// @Description Initiates Novofon callback that connects client number with user number; proxy_number is optional
// @Tags Calls
// @Accept json
// @Produce json
// @Param request body connectCallRequest true "Connect call payload"
// @Success 200 {object} calls.CallbackResponse "Call request accepted"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 502 {object} map[string]string "Provider error"
// @Failure 503 {object} map[string]string "Provider not configured"
// @Router /api/v1/calls/connect [post]
func (h *CallHandler) ConnectClientWithUserViaProxy(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		http.Error(w, calls.ErrProviderNotConfigured.Error(), http.StatusServiceUnavailable)
		return
	}

	var req connectCallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.service.ConnectClientWithUserViaProxy(r.Context(), calls.ConnectInput{
		ClientNumber: req.ClientNumber,
		UserNumber:   req.UserNumber,
		ProxyNumber:  req.ProxyNumber,
		Predicted:    req.Predicted,
	})
	if err != nil {
		switch {
		case errors.Is(err, calls.ErrClientNumberRequired),
			errors.Is(err, calls.ErrUserNumberRequired):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, calls.ErrProviderNotConfigured):
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		case errors.Is(err, novofon.ErrAPI):
			http.Error(w, err.Error(), http.StatusBadGateway)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ReceiveCallback godoc
// @Summary Receive Novofon callback webhook
// @Description Receives call events from Novofon PBX webhook and validates Signature header
// @Tags Calls
// @Accept x-www-form-urlencoded
// @Produce json
// @Param Signature header string false "Novofon signature"
// @Success 200 {object} map[string]interface{} "Callback accepted"
// @Failure 400 {object} map[string]string "Invalid callback payload"
// @Failure 401 {object} map[string]string "Invalid callback signature"
// @Failure 503 {object} map[string]string "Provider is not configured"
// @Router /api/v1/calls/callback [post]
func (h *CallHandler) ReceiveCallback(w http.ResponseWriter, r *http.Request) {
	if echo := r.URL.Query().Get("zd_echo"); echo != "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(echo))
		return
	}

	if h.service == nil {
		http.Error(w, calls.ErrProviderNotConfigured.Error(), http.StatusServiceUnavailable)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid callback payload", http.StatusBadRequest)
		return
	}

	signature := strings.TrimSpace(r.Header.Get("Signature"))
	event, err := h.service.HandleProviderCallback(r.PostForm, signature)
	if err != nil {
		switch {
		case errors.Is(err, calls.ErrCallbackNotSupported),
			errors.Is(err, calls.ErrProviderNotConfigured):
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		case errors.Is(err, novofon.ErrMissingWebhookSignature),
			errors.Is(err, novofon.ErrInvalidWebhookSignature):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case errors.Is(err, novofon.ErrUnsupportedWebhookEvent):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"event":  event,
	})
}
