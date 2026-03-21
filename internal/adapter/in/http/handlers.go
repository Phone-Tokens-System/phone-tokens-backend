package http

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/app"
)

type Handlers struct {
	User        *UserHandler
	Token       *TokenHandler
	Sms         *SmsHandler
	Agent       *AgentHandler
	Admin       *AdminHandler
	Billing     *BillingHandler
	UserProfile *UserProfileHandler
	Dict        *DictionaryHandler
}

func BuildHandlers(services app.Services) *Handlers {
	user := NewUserHandler(services.User)
	token := NewTokenHandler(services.Token)
	sms := NewSmsHandler(services.SMS)
	agent := NewAgentHandler(services.Cert, services.SMS)
	admin := NewAdminHandler(services.Cert)
	billing := NewBillingHandler(services.Billing)
	userProfile := NewUserProfileHandler(services.UserProfile)
	dict := DictionaryHandler{}
	return &Handlers{
		User:        user,
		Token:       token,
		Sms:         sms,
		Agent:       agent,
		Admin:       admin,
		Billing:     billing,
		UserProfile: userProfile,
		Dict:        &dict,
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
