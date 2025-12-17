package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"time"

	"phone-tokens/internal/model"
	"phone-tokens/internal/service/tokens"
)

type TokenHandler struct {
	service tokens.Service
}

func NewTokenHandler(service tokens.Service) *TokenHandler {
	return &TokenHandler{service: service}
}

type createTokenRequest struct {
	Name        string                  `json:"name"`
	Permissions []model.TokenPermission `json:"permissions"`
	TTLSeconds  int64                   `json:"ttl_seconds"`
}

type tokenResponse struct {
	ID          string                  `json:"id"`
	Token       string                  `json:"token"`
	Name        string                  `json:"name"`
	Permissions []model.TokenPermission `json:"permissions"`
	Status      model.TokenStatus       `json:"status"`
	ExpiresAt   string                  `json:"expires_at"`
	AgentId     string                  `json:"agent_id"`
}

// CreateToken godoc
// @Summary Issue a persistent token
// @Security BearerAuth
// @Description Issues a persistent token for the authenticated user
// @Tags Token
// @Accept json
// @Produce json
// @Param request body createTokenRequest true "Token request payload"
// @Success 201 {object} tokenResponse "Token successfully issued"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/tokens [post]
func (h *TokenHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
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

	token, err := h.service.Issue(r.Context(), tokens.IssueInput{
		UserID:      claims.UserID,
		Name:        req.Name,
		Permissions: req.Permissions,
		TTLSeconds:  req.TTLSeconds,
	})
	if err != nil {
		if isValidationError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, toTokenResponse(token))
}

type updateTokenRequest struct {
	TTLSeconds int64 `json:"ttl_seconds"`
}

// UpdateTokenTTL godoc
// @Summary Update token expiration
// @Security BearerAuth
// @Description Updates token expiration time for the authenticated user
// @Tags Token
// @Accept json
// @Produce json
// @Param tokenID path string true "Token ID"
// @Param request body updateTokenRequest true "Update TTL request payload"
// @Success 200 {object} tokenResponse "Token successfully updated"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Token not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/tokens/{tokenID} [patch]
func (h *TokenHandler) UpdateTokenTTL(w http.ResponseWriter, r *http.Request) {
	var req updateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.TTLSeconds <= 0 {
		http.Error(w, "ttl_seconds must be greater than zero", http.StatusBadRequest)
		return
	}

	tokenID := r.PathValue("tokenID")
	if tokenID == "" {
		http.Error(w, "token id is required", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := h.service.UpdateTTL(r.Context(), claims.UserID, tokenID, req.TTLSeconds)
	if err != nil {
		switch err {
		case tokens.ErrForbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case tokens.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(token))
}

// DeleteToken godoc
// @Summary Delete a token
// @Security BearerAuth
// @Description Removes a token for the authenticated user
// @Tags Token
// @Param tokenID path string true "Token ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Token not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/tokens/{tokenID} [delete]
func (h *TokenHandler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	tokenID := r.PathValue("tokenID")
	if tokenID == "" {
		http.Error(w, "token id is required", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.Delete(r.Context(), claims.UserID, tokenID)
	if err != nil {
		switch err {
		case tokens.ErrForbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case tokens.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FreezeToken godoc
// @Summary Freeze a token
// @Security BearerAuth
// @Description Sets token status to frozen
// @Tags Token
// @Param tokenID path string true "Token ID"
// @Success 200 {object} tokenResponse "Token successfully frozen"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Token not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/tokens/{tokenID}/freeze [patch]
func (h *TokenHandler) FreezeToken(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, model.TokenStatusFrozen)
}

// UnfreezeToken godoc
// @Summary Unfreeze a token
// @Security BearerAuth
// @Description Sets token status to active
// @Tags Token
// @Param tokenID path string true "Token ID"
// @Success 200 {object} tokenResponse "Token successfully unfrozen"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Token not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/tokens/{tokenID}/unfreeze [patch]
func (h *TokenHandler) UnfreezeToken(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, model.TokenStatusActive)
}

func (h *TokenHandler) changeStatus(w http.ResponseWriter, r *http.Request, status model.TokenStatus) {
	tokenID := r.PathValue("tokenID")
	if tokenID == "" {
		http.Error(w, "token id is required", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := h.service.SetStatus(r.Context(), claims.UserID, tokenID, status)
	if err != nil {
		switch err {
		case tokens.ErrForbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case tokens.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			if isValidationError(err) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, toTokenResponse(token))
}

func toTokenResponse(token *model.UserToken) tokenResponse {
	return tokenResponse{
		ID:          token.ID,
		Token:       token.Token,
		Name:        token.Name,
		Permissions: []model.TokenPermission(token.Permissions),
		Status:      token.Status,
		ExpiresAt:   token.ExpiresAt.Format(time.RFC3339),
		AgentId:     token.AgentId.String(),
	}
}

func isValidationError(err error) bool {
	return errors.Is(err, tokens.ErrInvalidPermission) || errors.Is(err, tokens.ErrInvalidStatus)
}

// BindAgentToToken godoc
// @Summary      Bind agent (service) to token
// @Description  Привязывает агента (внешний сервис) к пользовательскому токену по имени токена и ид агента
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        user_id path string true "User ID"
// @Success      200      {object}  tokenResponse
// @Failure      400      {string}  string  "invalid request body"
// @Failure      401      {string}  string  "unauthorized"
// @Failure      403      {string}  string  "forbidden"
// @Failure      404      {string}  string  "not found"
// @Failure      500      {string}  string  "internal server error"
// @Router       /tokens/bind-agent [post]
func (h *TokenHandler) BindAgentToToken(w http.ResponseWriter, r *http.Request) {
	var req dto.BindTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	updatedToken, err := h.service.BingAgentToTokenByName(r.Context(), claims.UserID, req)
	if err != nil {
		switch err {
		case tokens.ErrForbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case tokens.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(toTokenResponse(updatedToken)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetTokensByUser godoc
// @Summary      Get tokens by user
// @Description  Получение токенов пользователя
// @Tags         tokens
// @Accept       json
// @Produce      json
// @Param        request  body      dto.BindTokenRequest  true  "Bind agent to token request"
// @Success      200      {array}  tokenResponse
// @Failure      400      {string}  string  "invalid request body"
// @Failure      401      {string}  string  "unauthorized"
// @Failure      403      {string}  string  "forbidden"
// @Failure      404      {string}  string  "not found"
// @Failure      500      {string}  string  "internal server error"
// @Router        /users/{user_id}/tokens [get]
func (h *TokenHandler) GetTokensByUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("user_id")
	if userId == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.UserID != userId {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	tokensByUser, err := h.service.GetTokensByUser(r.Context(), userId)
	if err != nil {
		return
	}
	err = json.NewEncoder(w).Encode(tokensByUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
