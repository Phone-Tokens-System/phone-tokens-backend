package http

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/app"
)

type Handlers struct {
	User    *UserHandler
	Token   *TokenHandler
	Sms     *SmsHandler
	Agent   *AgentHandler
	Admin   *AdminHandler
	Billing *BillingHandler
}

func BuildHandlers(services app.Services) *Handlers {
	user := NewUserHandler(services.User)
	token := NewTokenHandler(services.Token)
	sms := NewSmsHandler(services.SMS)
	agent := NewAgentHandler(services.Cert)
	admin := NewAdminHandler(services.Cert)
	billing := NewBillingHandler(services.Billing)
	return &Handlers{
		User:    user,
		Token:   token,
		Sms:     sms,
		Agent:   agent,
		Admin:   admin,
		Billing: billing,
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
