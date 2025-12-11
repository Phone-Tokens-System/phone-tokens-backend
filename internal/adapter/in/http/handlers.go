package http

import (
	"encoding/json"
	"net/http"
	"phone-tokens/internal/app"
)

type Handlers struct {
	User  *UserHandler
	Token *TokenHandler
	Sms   *SmsHandler
	Agent *AgentHandler
}

func BuildHandlers(services app.Services) *Handlers {
	user := NewUserHandler(services.User)
	token := NewTokenHandler(services.Token)
	sms := NewSmsHandler(services.SMS)
	agent := NewAgentHandler(services.Cert)
	return &Handlers{
		User:  user,
		Token: token,
		Sms:   sms,
		Agent: agent,
	}
}

//type Handler struct {
//	UserHandler  *UserHandler
//	TokenHandler *TokenHandler
//	AgentHandler *AgentHandler
//	SmsHandler   *SmsHandler
//}
//
//func NewHandler(h app.Handlers) *Handler {
//	return &Handler{
//		h.User,
//		h.Token, h.Agent, h.Sms,
//	}
//}

//	type registerRequest struct {
//		Phone    string     `json:"phone"`
//		Password string     `json:"password"`
//		Role     model.Role `json:"role"`
//	}
//
//	type registerResponse struct {
//		ID    string     `json:"id"`
//		Phone string     `json:"phone"`
//		Role  model.Role `json:"role"`
//	}
//
//	func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
//		var req registerRequest
//		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//			http.Error(w, "invalid request body", http.StatusBadRequest)
//			return
//		}
//
//		user, err := h.userService.Register(r.Context(), req.Phone, req.Password, req.Role)
//		if err != nil {
//			switch err {
//			case users.ErrPhoneAlreadyUsed:
//				http.Error(w, err.Error(), http.StatusConflict)
//			default:
//				http.Error(w, "internal error", http.StatusInternalServerError)
//			}
//			return
//		}
//
//		resp := registerResponse{
//			ID:    user.ID,
//			Phone: user.Phone,
//			Role:  user.Role,
//		}
//
//		writeJSON(w, http.StatusCreated, resp)
//	}
//
//	type loginRequest struct {
//		Phone    string `json:"phone"`
//		Password string `json:"password"`
//	}
//
//	type loginResponse struct {
//		Token string `json:"token"`
//	}
//
//	func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
//		var req loginRequest
//		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//			http.Error(w, "invalid request body", http.StatusBadRequest)
//			return
//		}
//
//		token, _, err := h.userService.Authenticate(r.Context(), req.Phone, req.Password)
//		if err != nil {
//			switch err {
//			case users.ErrInvalidCredentials:
//				http.Error(w, err.Error(), http.StatusUnauthorized)
//			default:
//				http.Error(w, "internal error", http.StatusInternalServerError)
//			}
//			return
//		}
//
//		writeJSON(w, http.StatusOK, loginResponse{Token: token})
//	}
//
//	type createTokenRequest struct {
//		TTLSeconds int64 `json:"ttl_seconds"`
//	}
//
//	type tokenResponse struct {
//		ID        string `json:"id"`
//		Token     string `json:"token"`
//		ExpiresAt string `json:"expires_at"`
//	}
//
// // CreateToken issues a persistent token for the authenticated user.
//
//	func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
//		var req createTokenRequest
//		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//			http.Error(w, "invalid request body", http.StatusBadRequest)
//			return
//		}
//		if req.TTLSeconds <= 0 {
//			http.Error(w, "ttl_seconds must be greater than zero", http.StatusBadRequest)
//			return
//		}
//
//		claims, ok := r.Context().Value(userContextKey).(*UserClaims)
//		if !ok || claims.UserID == "" {
//			http.Error(w, "unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		token, err := h.tokenService.Issue(r.Context(), claims.UserID, req.TTLSeconds)
//		if err != nil {
//			http.Error(w, "internal error", http.StatusInternalServerError)
//			return
//		}
//
//		resp := tokenResponse{
//			ID:        token.ID,
//			Token:     token.Token,
//			ExpiresAt: token.ExpiresAt.Format(time.RFC3339),
//		}
//
//		writeJSON(w, http.StatusCreated, resp)
//	}
//
//	type updateTokenRequest struct {
//		TTLSeconds int64 `json:"ttl_seconds"`
//	}
//
// // UpdateTokenTTL updates token expiration time for the authenticated user.
//
//	func (h *Handler) UpdateTokenTTL(w http.ResponseWriter, r *http.Request) {
//		var req updateTokenRequest
//		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//			http.Error(w, "invalid request body", http.StatusBadRequest)
//			return
//		}
//		if req.TTLSeconds <= 0 {
//			http.Error(w, "ttl_seconds must be greater than zero", http.StatusBadRequest)
//			return
//		}
//
//		tokenID := r.PathValue("tokenID")
//		if tokenID == "" {
//			http.Error(w, "token id is required", http.StatusBadRequest)
//			return
//		}
//
//		claims, ok := r.Context().Value(userContextKey).(*UserClaims)
//		if !ok || claims.UserID == "" {
//			http.Error(w, "unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		token, err := h.tokenService.UpdateTTL(r.Context(), claims.UserID, tokenID, req.TTLSeconds)
//		if err != nil {
//			switch err {
//			case tokens.ErrForbidden:
//				http.Error(w, err.Error(), http.StatusForbidden)
//			case tokens.ErrNotFound:
//				http.Error(w, err.Error(), http.StatusNotFound)
//			default:
//				http.Error(w, "internal error", http.StatusInternalServerError)
//			}
//			return
//		}
//
//		resp := tokenResponse{
//			ID:        token.ID,
//			Token:     token.Token,
//			ExpiresAt: token.ExpiresAt.Format(time.RFC3339),
//		}
//
//		writeJSON(w, http.StatusOK, resp)
//	}
//
// // DeleteToken removes a token for the authenticated user.
//
//	func (h *Handler) DeleteToken(w http.ResponseWriter, r *http.Request) {
//		tokenID := r.PathValue("tokenID")
//		if tokenID == "" {
//			http.Error(w, "token id is required", http.StatusBadRequest)
//			return
//		}
//
//		claims, ok := r.Context().Value(userContextKey).(*UserClaims)
//		if !ok || claims.UserID == "" {
//			http.Error(w, "unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		err := h.tokenService.Delete(r.Context(), claims.UserID, tokenID)
//		if err != nil {
//			switch err {
//			case tokens.ErrForbidden:
//				http.Error(w, err.Error(), http.StatusForbidden)
//			case tokens.ErrNotFound:
//				http.Error(w, err.Error(), http.StatusNotFound)
//			default:
//				http.Error(w, "internal error", http.StatusInternalServerError)
//			}
//			return
//		}
//
//		w.WriteHeader(http.StatusNoContent)
//	}
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
