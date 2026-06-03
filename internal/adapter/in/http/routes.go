package http

import (
	"net/http"

	_ "phone-tokens/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(mux *http.ServeMux, h Handlers, authCfg AuthConfig) {
	mux.HandleFunc("POST /api/v1/register", h.User.Register)
	mux.HandleFunc("POST /api/v1/login", h.User.Login)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	authMiddleware := AuthMiddleware(authCfg)

	mux.Handle("/api/v1/me", authMiddleware(http.HandlerFunc(h.User.Me)))

	mux.Handle("/api/v1/admin/ping", authMiddleware(RequireRole("admin", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("admin ok"))
	}))))

	// Issue user tokens (requires authentication).
	mux.Handle("/api/v1/tokens", authMiddleware(http.HandlerFunc(h.Token.CreateToken)))
	mux.Handle("PATCH /api/v1/tokens/{tokenID}", authMiddleware(http.HandlerFunc(h.Token.UpdateTokenTTL)))
	mux.Handle("DELETE /api/v1/tokens/{tokenID}", authMiddleware(http.HandlerFunc(h.Token.DeleteToken)))
	mux.Handle("PATCH /api/v1/tokens/{tokenID}/freeze", authMiddleware(http.HandlerFunc(h.Token.FreezeToken)))
	mux.Handle("PATCH /api/v1/tokens/{tokenID}/unfreeze", authMiddleware(http.HandlerFunc(h.Token.UnfreezeToken)))
	mux.Handle("GET /api/v1/users/{userId}/tokens", authMiddleware(http.HandlerFunc(h.Token.GetTokensByUser)))
	mux.Handle("POST /api/v1/tokens/bind-agent", authMiddleware(http.HandlerFunc(h.Token.BindAgentToToken)))

	//manage user profiles (user)
	mux.Handle("GET /api/v1/user-profile/filters", http.HandlerFunc(h.UserProfile.GetFilters))
	mux.Handle("GET /api/v1/userprofile/filters", http.HandlerFunc(h.UserProfile.GetFilters))
	mux.Handle("POST /api/v1/user-profile", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.SaveUserProfile))))
	mux.Handle("DELETE /api/v1/user-profile", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.DeleteUserProfile))))
	mux.Handle("PUT /api/v1/user-profile", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.UpdateUserProfile))))
	mux.Handle("GET /api/v1/user-profile/me", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.GetUserProfileById))))

	// agents
	mux.Handle("GET /api/v1/csr/signed/current", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Agent.GetCurrentSignedCertificate))))
	mux.Handle("GET /api/v1/csr/signed", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Agent.GetSignedCertificate))))
	mux.Handle("POST /api/v1/csr/upload", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Agent.UploadCSR))))
	mux.Handle("GET /api/v1/sms/agents/logs", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Agent.SeeSMSLogs))))
	// manage user profiles(agent)
	mux.Handle("POST /api/v1/agents/tokens/user-profile", authMiddleware(RequireRole("agent", http.HandlerFunc(h.UserProfile.GetUserProfileByToken))))
	mux.Handle("GET /api/v1/agents/tokens/user-profile", authMiddleware(RequireRole("agent", http.HandlerFunc(h.UserProfile.GetUserProfilesByAgentID))))
	mux.Handle("POST /api/v1/agents/tokens/user-profile/filtered", authMiddleware(RequireRole("agent", http.HandlerFunc(h.UserProfile.GetUserProfilesFilteredByAgentID))))

	//admin
	mux.Handle("POST /api/v1/admin/csr", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Admin.AcceptCSRRequest))))
	mux.Handle("GET /api/v1/admin/csr", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Admin.ShowCSRRequests))))
	mux.Handle("POST /api/v1/admin/csr/approve/{id}", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Admin.ApproveCSRRequest))))

	//sms
	mux.Handle("POST /api/v1/sms/send", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Sms.sendSMS))))
	mux.Handle("POST /api/v1/sms/send_filtered", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Sms.SendSmsWithFilters))))
	mux.Handle("GET /api/v1/sms/agents/{agentId}", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Sms.getSmsListByAgentId))))

	// sms admin
	mux.Handle("GET /api/v1/sms/logs", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.getSmsList))))
	mux.Handle("GET /api/v1/sms/status", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.checkStatus))))
	mux.Handle("GET /api/v1/sms/all", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.getSmsListFromProvider))))

	mux.Handle("GET /api/v1/sms/users/{token}", authMiddleware(http.HandlerFunc(h.Sms.getSmsListByToken)))

	// billing
	mux.Handle("POST /api/v1/billing/balance", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Billing.TopBalance))))
	mux.HandleFunc("POST /api/v1/billing/webhook", h.Billing.StripeWebhookHandler)
	mux.Handle("GET /api/v1/billing/{agent_id}/balance", authMiddleware(http.HandlerFunc(h.Billing.GetBalanceHandler)))
	mux.Handle("GET /api/v1/billing/balance", authMiddleware(http.HandlerFunc(h.Billing.GetBalanceHandler)))

	// accounting (agent transaction history)
	mux.Handle("GET /api/v1/agents/{agent_id}/transactions", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Billing.GetTransactionsHandler))))

	// packages
	mux.Handle("GET /api/v1/packages", authMiddleware(http.HandlerFunc(h.Billing.SeePackageOptionsHandler)))
	mux.Handle("POST /api/v1/agents/{agent_id}/packages", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Billing.BuyPackageHandler))))
	mux.Handle("GET /api/v1/agents/{agent_id}/packages", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Billing.SeeAgentPackagesHandler))))

	// admin: управление пакетами
	mux.Handle("POST /api/v1/admin/packages", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Billing.CreatePackageHandler))))
	mux.Handle("DELETE /api/v1/admin/packages/{pkg_id}", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Billing.DeletePackageHandler))))

	// SSO — привязка токена пользователя к внешнему агенту
	// GET  /api/v1/sso/authorize  — агент перенаправляет сюда пользователя (без auth)
	// POST /api/v1/sso/complete   — frontend вызывает после логина пользователя (требует user JWT)
	// GET  /api/v1/sso/me         — агент проверяет валидность полученного токена (без auth)
	mux.HandleFunc("GET /api/v1/sso/authorize", h.SSO.Authorize)
	mux.Handle("POST /api/v1/sso/complete", authMiddleware(RequireRole("user", http.HandlerFunc(h.SSO.Complete))))
	mux.HandleFunc("GET /api/v1/sso/me", h.SSO.Me)

	//dict
	mux.Handle("GET /api/v1/dictionary/countries", http.HandlerFunc(h.Dict.GetCountries))
	mux.Handle("GET /api/v1/dictionary/regions", http.HandlerFunc(h.Dict.GetRegions))
	mux.Handle("GET /api/v1/dictionary/cities", http.HandlerFunc(h.Dict.GetCities))
}
