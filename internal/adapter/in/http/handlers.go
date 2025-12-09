package http

import (
	"encoding/json"
	"net/http"

	"phone_token_system/internal/service/tokens"
	"phone_token_system/internal/service/users"
)

type Handler struct {
	Users  *UserHandler
	Tokens *TokenHandler
}

func NewHandler(userSvc users.Service, tokenSvc tokens.Service) *Handler {
	return &Handler{
		Users:  NewUserHandler(userSvc),
		Tokens: NewTokenHandler(tokenSvc),
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
