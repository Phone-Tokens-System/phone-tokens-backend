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

	mux.HandleFunc("POST /api/v1/csr", h.Agent.AcceptCSRRequest)
	mux.HandleFunc("GET /api/v1/csr/signed", h.Agent.GetSignedCertificate)
	mux.HandleFunc("GET /api/v1/admin/csr", h.Agent.ShowCSRRequests)
	mux.HandleFunc("POST /api/v1/admin/csr/approve/{id}", h.Agent.ApproveCSRRequest)
	mux.HandleFunc("POST /api/v1/csr/upload", h.Agent.UploadCSR)
	mux.HandleFunc("GET /api/v1/sms/list", h.Sms.getSmsList)
	mux.HandleFunc("POST /api/v1/sms/send", h.Sms.sendSMS)
	mux.HandleFunc("GET /api/v1/sms/status", h.Sms.checkStatus)
	mux.HandleFunc("GET /api/v1/sms/all", h.Sms.getSmsListFromProvider)
}
