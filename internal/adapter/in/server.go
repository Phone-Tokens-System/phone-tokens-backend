package in

import (
	"errors"
	"net/http"
	httpadapter "phone-tokens/internal/adapter/in/http"
)

// NewHTTPServer создаёт HTTP‑сервер поверх доменного сервиса.
func NewHTTPServer(httpPort string, jwtSecret string, handlers httpadapter.Handlers) (*http.Server, error) {
	if httpPort == "" {
		return nil, errors.New("HTTP_PORT is required")
	}

	mux := http.NewServeMux()
	httpadapter.RegisterRoutes(mux, handlers, httpadapter.AuthConfig{
		JWTSecret: jwtSecret,
	})

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	return server, nil
}
