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

	mux.Handle("/api/v1/me", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userContextKey).(*UserClaims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		writeJSON(w, http.StatusOK, claims)
	})))

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
	mux.Handle("POST /api/v1/user-profile", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.SaveUserProfile))))
	mux.Handle("DELETE /api/v1/user-profile", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.DeleteUserProfile))))
	mux.Handle("PUT /api/v1/user-profile", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.UpdateUserProfile))))
	mux.Handle("GET /api/v1/user-profile/me", authMiddleware(RequireRole("user", http.HandlerFunc(h.UserProfile.GetUserProfileById))))

	// agents
	mux.HandleFunc("GET /api/v1/csr/signed", h.Agent.GetSignedCertificate)
	mux.HandleFunc("POST /api/v1/csr/upload", h.Agent.UploadCSR)
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
	mux.Handle("POST /api/v1/sms/send_filtered", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Sms.sendSMS))))
	mux.Handle("GET /api/v1/sms/agents/{agentId}", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Sms.getSmsListByAgentId))))

	// sms admin
	mux.Handle("GET /api/v1/sms/logs", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.getSmsList))))
	mux.Handle("GET /api/v1/sms/status", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.checkStatus))))
	mux.Handle("GET /api/v1/sms/all", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.getSmsListFromProvider))))

	mux.Handle("GET /api/v1/sms/users/{token}", authMiddleware(http.HandlerFunc(h.Sms.getSmsListByToken)))

	// billing
	mux.Handle("POST /api/v1/billing/balance", authMiddleware(RequireRole("agent", http.HandlerFunc(h.Billing.TopBalance))))
	mux.HandleFunc("POST /api/v1/billing/webhook", h.Billing.StripeWebhookHandler)
	mux.Handle("GET /api/v1/billing/balance", authMiddleware(http.HandlerFunc(h.Billing.GetBalanceHandler)))
}
