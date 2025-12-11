package http

import (
	"encoding/json"
	"errors"
	"net/http"
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
}

// CreateToken issues a persistent token for the authenticated user.
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

// UpdateTokenTTL updates token expiration time for the authenticated user.
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

// DeleteToken removes a token for the authenticated user.
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

func (h *TokenHandler) FreezeToken(w http.ResponseWriter, r *http.Request) {
	h.changeStatus(w, r, model.TokenStatusFrozen)
}

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
		Permissions: token.Permissions,
		Status:      token.Status,
		ExpiresAt:   token.ExpiresAt.Format(time.RFC3339),
	}
}

func isValidationError(err error) bool {
	return errors.Is(err, tokens.ErrInvalidPermission) || errors.Is(err, tokens.ErrInvalidStatus)
}
