package http

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handler, authCfg AuthConfig) {
	mux.HandleFunc("POST /api/v1/register", h.Register)
	mux.HandleFunc("POST /api/v1/login", h.Login)

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
}
