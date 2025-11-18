package http

import (
	"encoding/json"
	"net/http"

	"users/internal/users/model"
	"users/internal/users/service/users"
)

type Handler struct {
	service users.Service
}

func NewHandler(s users.Service) *Handler {
	return &Handler{service: s}
}

type registerRequest struct {
	Phone    string      `json:"phone"`
	Password string      `json:"password"`
	Role     model.Role  `json:"role"`
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

	user, err := h.service.Register(r.Context(), req.Phone, req.Password, req.Role)
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

	token, _, err := h.service.Authenticate(r.Context(), req.Phone, req.Password)
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

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

