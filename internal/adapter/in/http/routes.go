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

	// certificates
	mux.HandleFunc("POST /api/v1/csr", h.Agent.AcceptCSRRequest)
	mux.HandleFunc("GET /api/v1/csr/signed", h.Agent.GetSignedCertificate)
	mux.Handle("GET /api/v1/admin/csr", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Agent.ShowCSRRequests))))
	mux.Handle("POST /api/v1/admin/csr/approve/{id}", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Agent.ApproveCSRRequest))))
	mux.HandleFunc("POST /api/v1/csr/upload", h.Agent.UploadCSR)

	//sms
	mux.Handle("GET /api/v1/sms/logs", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.getSmsList))))
	mux.HandleFunc("POST /api/v1/sms/send", h.Sms.sendSMS)
	mux.HandleFunc("GET /api/v1/sms/status", h.Sms.checkStatus)
	mux.Handle("GET /api/v1/sms/all", authMiddleware(RequireRole("admin", http.HandlerFunc(h.Sms.getSmsListFromProvider))))

	mux.Handle("GET /api/v1/sms/users/userId", authMiddleware(http.HandlerFunc(h.Sms.getSmsListByUser)))
	mux.Handle("GET /api/v1/sms/agents/agentId", authMiddleware(http.HandlerFunc(h.Sms.getSmsListByAgentId)))
}
