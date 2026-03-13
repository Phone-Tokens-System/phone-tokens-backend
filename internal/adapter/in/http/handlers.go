package http

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/app"
)

type Handlers struct {
	User  *UserHandler
	Token *TokenHandler
	Call  *CallHandler
	Sms   *SmsHandler
	Agent *AgentHandler
}

func BuildHandlers(services app.Services) *Handlers {
	user := NewUserHandler(services.User)
	token := NewTokenHandler(services.Token)
	call := NewCallHandler(services.Call)
	sms := NewSmsHandler(services.SMS)
	agent := NewAgentHandler(services.Cert)
	return &Handlers{
		User:  user,
		Token: token,
		Call:  call,
		Sms:   sms,
		Agent: agent,
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
