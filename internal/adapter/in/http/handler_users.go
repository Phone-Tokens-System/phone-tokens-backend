package http

import (
	"encoding/json"
	"errors"
	"fmt"
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
// @Param request body users.RegisterRequest true "User registration payload"
// @Success 201 {object} registerResponse "User successfully registered"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 409 {object} map[string]string "Phone already used"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req users.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		fmt.Println(req)
		return
	}
	fmt.Println(req)
	user, err := h.service.Register(r.Context(), req)
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

type meResponse struct {
	UserID  string     `json:"user_id"`
	Phone   string     `json:"phone"`
	Role    model.Role `json:"role"`
	AgentID string     `json:"agent_id,omitempty"`
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
		fmt.Println(req)
		return
	}
	fmt.Println(req)
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

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	resp := meResponse{
		UserID: claims.UserID,
		Phone:  claims.Phone,
		Role:   claims.Role,
	}

	if claims.Role == model.RoleAgent {
		agent, err := h.service.GetAgentByUserID(r.Context(), claims.UserID)
		if err != nil {
			if !errors.Is(err, model.ErrNotFound) {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		} else if agent != nil {
			resp.AgentID = agent.ID
		}
	}

	writeJSON(w, http.StatusOK, resp)
}
