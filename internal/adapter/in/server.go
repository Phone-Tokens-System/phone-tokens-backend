package in

import (
	"errors"
	"net/http"

	httpadapter "phone_token_system/internal/adapter/in/http"
	"phone_token_system/internal/service/tokens"
	"phone_token_system/internal/service/users"
)

// NewHTTPServer создаёт HTTP‑сервер поверх доменного сервиса.
func NewHTTPServer(httpPort string, jwtSecret string, userSvc users.Service, tokenSvc tokens.Service) (*http.Server, error) {
	if httpPort == "" {
		return nil, errors.New("HTTP_PORT is required")
	}

	handler := httpadapter.NewHandler(userSvc, tokenSvc)

	mux := http.NewServeMux()
	httpadapter.RegisterRoutes(mux, handler, httpadapter.AuthConfig{
		JWTSecret: jwtSecret,
	})

	server := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	return server, nil
}
