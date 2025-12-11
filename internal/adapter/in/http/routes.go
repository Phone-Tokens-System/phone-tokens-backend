package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler, authCfg AuthConfig) {
	mux.HandleFunc("POST /api/v1/register", h.Users.Register)
	mux.HandleFunc("POST /api/v1/login", h.Users.Login)

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
	mux.Handle("/api/v1/tokens", authMiddleware(http.HandlerFunc(h.Tokens.CreateToken)))
	mux.Handle("PATCH /api/v1/tokens/{tokenID}", authMiddleware(http.HandlerFunc(h.Tokens.UpdateTokenTTL)))
	mux.Handle("DELETE /api/v1/tokens/{tokenID}", authMiddleware(http.HandlerFunc(h.Tokens.DeleteToken)))
	mux.Handle("PATCH /api/v1/tokens/{tokenID}/freeze", authMiddleware(http.HandlerFunc(h.Tokens.FreezeToken)))
	mux.Handle("PATCH /api/v1/tokens/{tokenID}/unfreeze", authMiddleware(http.HandlerFunc(h.Tokens.UnfreezeToken)))
}
