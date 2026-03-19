package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signedToken(t *testing.T, secret, role string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub":   "user-1",
		"phone": "79990000000",
		"role":  role,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signed
}

func makeHandlers() Handlers {
	return Handlers{
		User:    &UserHandler{},
		Token:   &TokenHandler{},
		Sms:     &SmsHandler{},
		Agent:   &AgentHandler{},
		Admin:   &AdminHandler{},
		Billing: &BillingHandler{},
	}
}

func TestSMSSendRequiresAgentRole_UserGetsForbidden(t *testing.T) {
	secret := "test-secret"

	mux := http.NewServeMux()
	RegisterRoutes(mux, makeHandlers(), AuthConfig{JWTSecret: secret})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sms/send", strings.NewReader(`{`))
	req.Header.Set("Authorization", "Bearer "+signedToken(t, secret, "user"))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestSMSSendAllowsAgentRoleToReachHandler(t *testing.T) {
	secret := "test-secret"

	mux := http.NewServeMux()
	RegisterRoutes(mux, makeHandlers(), AuthConfig{JWTSecret: secret})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sms/send", strings.NewReader(`{`))
	req.Header.Set("Authorization", "Bearer "+signedToken(t, secret, "agent"))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// role guard passed; handler returns 400 due invalid JSON
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
