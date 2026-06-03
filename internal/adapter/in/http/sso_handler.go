package http

// SSO (Single Sign-On) flow — привязка токена пользователя к внешнему агенту.
//
// Сценарий:
//   1. Внешний сервис (агент) перенаправляет пользователя на наш endpoint:
//        GET /api/v1/sso/authorize?agent_id=<UUID>&redirect_uri=<callback>&state=<opaque>
//   2. Мы валидируем agent_id и перенаправляем пользователя на frontend нашей системы.
//        Там пользователь логинится (если не залогинен) и подтверждает привязку токена.
//   3. Frontend вызывает нас:
//        POST /api/v1/sso/complete   (requires Bearer JWT)
//        Body: { "agent_id", "redirect_uri", "state", "token_name", "permissions" }
//   4. Мы выдаём/находим UserToken для пары (user, agent) и отвечаем JSON:
//        { "redirect_uri": "<callback>?token=<token>&state=<state>" }
//   5. Frontend редиректит пользователя по этому URL.
//   6. Внешний сервис получает token и использует его как client_token при вызовах API
//      (POST /api/v1/sms/send и т.д.).

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/tokens"
	"phone-tokens/internal/service/users"
	"time"
)

type SSOHandler struct {
	users       users.Service
	tokenSvc    tokens.Service
	frontendURL string // e.g. "http://localhost:5173"
}

func NewSSOHandler(u users.Service, t tokens.Service, frontendURL string) *SSOHandler {
	return &SSOHandler{users: u, tokenSvc: t, frontendURL: frontendURL}
}

// Authorize godoc
// @Summary SSO: перенаправить пользователя на страницу логина
// @Description Внешний агент перенаправляет сюда пользователя. Мы проверяем agent_id
//
//	и отправляем пользователя на наш frontend для авторизации.
//
// @Tags SSO
// @Param agent_id    query string true  "UUID зарегистрированного агента"
// @Param redirect_uri query string true  "URL возврата внешнего сервиса"
// @Param state        query string false "Произвольное значение для защиты от CSRF"
// @Success 302 "Redirect to frontend login"
// @Failure 400 {string} string "Bad request"
// @Router /api/v1/sso/authorize [get]
func (h *SSOHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	if agentID == "" || redirectURI == "" {
		http.Error(w, "agent_id and redirect_uri are required", http.StatusBadRequest)
		return
	}

	// Проверяем что агент существует в системе
	_, err := h.users.GetAgentByID(r.Context(), agentID)
	if err != nil {
		http.Error(w, "unknown agent", http.StatusBadRequest)
		return
	}

	// Перенаправляем на страницу SSO нашего frontend
	// Там пользователь логинится и подтверждает выдачу токена агенту
	frontendLogin := fmt.Sprintf(
		"%s/sso?agent_id=%s&redirect_uri=%s&state=%s",
		h.frontendURL,
		url.QueryEscape(agentID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)
	log.Printf("frontendURL=%q", h.frontendURL)
	log.Printf("frontendLogin=%q", frontendLogin)

	http.Redirect(w, r, frontendLogin, http.StatusFound)
}

type SSOCompleteRequest struct {
	AgentID     string   `json:"agent_id"`
	RedirectURI string   `json:"redirect_uri"`
	State       string   `json:"state"`
	TokenName   string   `json:"token_name"`            // необязательно, для названия токена
	Permissions []string `json:"permissions,omitempty"` // ["sms", "calls"]
	TTLDays     int      `json:"ttl_days,omitempty"`    // 0 = бессрочный (1 год)
}

type SSOCompleteResponse struct {
	// Готовый URL для редиректа пользователя обратно к агенту.
	// Frontend должен сделать window.location.href = redirect_url
	RedirectURL string `json:"redirect_url"`
	// Токен, который агент получит в query-параметре ?token=
	Token string `json:"token"`
}

// Complete godoc
// @Summary SSO: выдать токен пользователя агенту
// @Description Вызывается нашим frontend после того как пользователь залогинился и подтвердил привязку.
//
//	Требует JWT-авторизации (пользователь).
//	Возвращает redirect_url — frontend делает window.location.href на него.
//
// @Tags SSO
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body SSOCompleteRequest true "Параметры SSO"
// @Success 200 {object} SSOCompleteResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/sso/complete [post]
func (h *SSOHandler) Complete(w http.ResponseWriter, r *http.Request) {
	// Требуем аутентификации пользователя
	claims, ok := r.Context().Value(userContextKey).(*UserClaims)
	if !ok || claims.UserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req SSOCompleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.AgentID == "" || req.RedirectURI == "" {
		http.Error(w, "agent_id and redirect_uri are required", http.StatusBadRequest)
		return
	}

	// Проверяем, что агент существует
	_, err := h.users.GetAgentByID(r.Context(), req.AgentID)
	if err != nil {
		http.Error(w, "unknown agent", http.StatusBadRequest)
		return
	}

	// Формируем разрешения токена
	perms := model.TokenPermissions{}
	for _, p := range req.Permissions {
		perms = append(perms, model.TokenPermission(p))
	}
	if len(perms) == 0 {
		perms = model.TokenPermissions{model.TokenPermissionSMS}
	}

	// Срок жизни токена
	ttlDays := req.TTLDays
	if ttlDays <= 0 {
		ttlDays = 365 // 1 год по умолчанию
	}
	ttlSeconds := int64(ttlDays * 24 * 3600)

	tokenName := req.TokenName
	if tokenName == "" {
		tokenName = "sso-token"
	}

	// Выдаём новый UserToken для этой пары (user, agent)
	issued, err := h.tokenSvc.Issue(r.Context(), tokens.IssueInput{
		UserID:      claims.UserID,
		AgentId:     req.AgentID,
		Name:        tokenName,
		TTLSeconds:  ttlSeconds,
		Permissions: perms,
	})
	if err != nil {
		http.Error(w, "failed to issue token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Собираем URL редиректа назад к агенту
	// ?token=<token_value>&state=<state>&expires_at=<unix>
	callbackURL, err := url.Parse(req.RedirectURI)
	if err != nil {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}
	q := callbackURL.Query()
	q.Set("token", issued.Token)
	if req.State != "" {
		q.Set("state", req.State)
	}
	q.Set("expires_at", fmt.Sprintf("%d", issued.ExpiresAt.Unix()))
	callbackURL.RawQuery = q.Encode()

	writeJSON(w, http.StatusOK, SSOCompleteResponse{
		RedirectURL: callbackURL.String(),
		Token:       issued.Token,
	})
}

// Me godoc (SSO)
// @Summary SSO: проверить валидность токена
// @Description Внешний агент проверяет, что токен, полученный после SSO, ещё действителен.
//
//	Реальный номер телефона пользователя НЕ раскрывается — только признак валидности.
//
// @Tags SSO
// @Param token query string true "Токен, полученный после SSO"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Router /api/v1/sso/me [get]
func (h *SSOHandler) Me(w http.ResponseWriter, r *http.Request) {
	tokenVal := r.URL.Query().Get("token")
	if tokenVal == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	// Проверяем существование токена (и его срок действия через репозиторий)
	// GetUserNumberFromToken внутри обращается к БД — если токена нет, вернёт ошибку.
	// Реальный номер телефона мы НЕ возвращаем агенту.
	_, err := h.tokenSvc.GetUserNumberFromToken(r.Context(), tokenVal)
	if err != nil {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"valid":      true,
		"checked_at": time.Now().UTC(),
	})
}
