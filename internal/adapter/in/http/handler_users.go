package http

import (
	"encoding/json"
	"net/http"

	"phone-tokens/internal/model"
	"phone-tokens/internal/service/users"
)

type UserHandler struct {
	service users.Service
}

func NewUserHandler(service users.Service) *UserHandler {
	return &UserHandler{service: service}
}

type registerRequest struct {
	Phone       string     `json:"phone"`
	Password    string     `json:"password"`
	Role        model.Role `json:"role"`
	ServiceName string     `json:"service_name"`
	Email       string     `json:"email"`
}

type registerResponse struct {
	ID    string     `json:"id"`
	Phone string     `json:"phone"`
	Role  model.Role `json:"role"`
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with phone, password and role
// @Tags User
// @Accept json
// @Produce json
// @Param request body registerRequest true "User registration payload"
// @Success 201 {object} registerResponse "User successfully registered"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 409 {object} map[string]string "Phone already used"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req.Phone, req.Password, req.Role, req.ServiceName, req.Email)
	if err != nil {
		switch err {
		case users.ErrPhoneAlreadyUsed:
			http.Error(w, err.Error(), http.StatusConflict)
		case users.ErrRoleNotAllowed:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case users.ErrAgentDetailsNeeded:
			http.Error(w, err.Error(), http.StatusBadRequest)
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

// Login godoc
// @Summary Login a user
// @Description Authenticates user with phone and password and returns a JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param request body loginRequest true "User login payload"
// @Success 200 {object} loginResponse "JWT token"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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
