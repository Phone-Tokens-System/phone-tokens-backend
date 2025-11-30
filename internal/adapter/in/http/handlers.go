package http

import (
	"encoding/json"
	"net/http"
	"time"

	"phone_token_system/internal/model"
	"phone_token_system/internal/service/tokens"
	"phone_token_system/internal/service/users"
)

type Handler struct {
	userService  users.Service
	tokenService tokens.Service
}

func NewHandler(userSvc users.Service, tokenSvc tokens.Service) *Handler {
	return &Handler{
		userService:  userSvc,
		tokenService: tokenSvc,
	}
}

type registerRequest struct {
	Phone    string     `json:"phone"`
	Password string     `json:"password"`
	Role     model.Role `json:"role"`
}

type registerResponse struct {
	ID    string     `json:"id"`
	Phone string     `json:"phone"`
	Role  model.Role `json:"role"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), req.Phone, req.Password, req.Role)
	if err != nil {
		switch err {
		case users.ErrPhoneAlreadyUsed:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	resp := registerResponse{
		ID:    user.ID,
		Phone: user.Phone,
		Role:  user.Role,
	}

	writeJSON(w, http.StatusCreated, resp)
}

type loginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, _, err := h.userService.Authenticate(r.Context(), req.Phone, req.Password)
	if err != nil {
		switch err {
		case users.ErrInvalidCredentials:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, loginResponse{Token: token})
}

type createTokenRequest struct {
	TTLSeconds int64 `json:"ttl_seconds"`
}

type createTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// CreateToken issues a persistent token for the authenticated user.
func (h *Handler) CreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.TTLSeconds <= 0 {
		http.Error(w, "ttl_seconds must be greater than zero", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := h.tokenService.Issue(r.Context(), claims.UserID, req.TTLSeconds)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := createTokenResponse{
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt.Format(time.RFC3339),
	}

	writeJSON(w, http.StatusCreated, resp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
