// Пример простого внешнего сервиса, который использует Phone Tokens SSO
// для привязки пользователя (без знания его реального номера телефона).
//
// Запуск:
//   go run examples/external_service/main.go
//
// Откройте http://localhost:9090 в браузере.
//
// Что делает этот сервис:
//   1. Показывает кнопку "Зарегистрироваться через Phone Tokens"
//   2. Перенаправляет пользователя на SSO нашей системы:
//        GET http://localhost:8080/api/v1/sso/authorize?agent_id=<AGENT_ID>&redirect_uri=http://localhost:9090/callback&state=<random>
//   3. Пользователь логинится на нашем frontend (http://localhost:5173/sso?...)
//      и подтверждает выдачу токена.
//   4. Наш frontend вызывает POST /api/v1/sso/complete и получает redirect_url.
//   5. Пользователь возвращается на http://localhost:9090/callback?token=<token>&state=<state>
//   6. Сервис сохраняет токен и может слать SMS через нашу систему.
//
// Переменные окружения:
//   AGENT_ID         — UUID зарегистрированного агента (из нашей системы)
//   AGENT_CERT_PATH  — путь к сертификату агента (PEM, для API-запросов)
//   PHONE_TOKENS_API — базовый URL API (default: http://localhost:8080)

package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// -------- конфигурация --------

//go:embed templates/*.html
var templateFS embed.FS

type Config struct {
	AgentID         string
	PhoneTokensAPI  string
	ListenAddr      string
	CallbackBaseURL string
}

func loadConfig() Config {
	return Config{
		AgentID:         getenv("AGENT_ID", "..."),
		PhoneTokensAPI:  getenv("PHONE_TOKENS_API", "http://localhost:8080"),
		ListenAddr:      getenv("LISTEN_ADDR", ":9090"),
		CallbackBaseURL: getenv("CALLBACK_BASE_URL", "http://localhost:9090"),
	}
}

var cfg Config

var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env file %v", err)
	}
	return err
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// -------- in-memory хранилище "пользователей" --------

type AppUser struct {
	Token    string // токен от нашей системы
	LinkedAt time.Time
}

var (
	mu    sync.Mutex
	users = map[string]*AppUser{} // state -> AppUser (до привязки state хранит ожидание)
)

// -------- обработчики --------

// GET / — главная страница с кнопкой регистрации
func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GET /start-sso — генерируем state и перенаправляем на наш SSO authorize endpoint
func startSSOHandler(w http.ResponseWriter, r *http.Request) {
	// Генерируем случайный state для защиты от CSRF
	stateBytes := make([]byte, 16)
	_, _ = rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	// Сохраняем ожидающий state
	mu.Lock()
	users[state] = nil // nil означает "ещё не привязан"
	mu.Unlock()

	// Строим URL редиректа на наш SSO
	redirectURI := cfg.CallbackBaseURL + "/callback"
	ssoURL := fmt.Sprintf(
		"%s/api/v1/sso/authorize?agent_id=%s&redirect_uri=%s&state=%s",
		cfg.PhoneTokensAPI,
		url.QueryEscape(cfg.AgentID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)

	log.Printf("SSO start: redirecting to %s", ssoURL)
	http.Redirect(w, r, ssoURL, http.StatusFound)
}

// GET /callback?token=<token>&state=<state>&expires_at=<unix>
// — возврат от нашего SSO после того как пользователь подтвердил привязку
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	state := r.URL.Query().Get("state")

	if token == "" || state == "" {
		http.Error(w, "missing token or state", http.StatusBadRequest)
		return
	}

	// Проверяем что state нам знаком (защита от CSRF)
	mu.Lock()
	_, known := users[state]
	mu.Unlock()
	if !known {
		http.Error(w, "unknown state — possible CSRF attempt", http.StatusBadRequest)
		return
	}

	// Опционально: верифицируем токен у нашего бэкенда
	valid, err := verifyToken(token)
	if err != nil || !valid {
		if err == nil {
			err = errors.New("invalid token")
		}
		http.Error(w, "token validation failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Сохраняем пользователя
	mu.Lock()
	users[state] = &AppUser{
		Token:    token,
		LinkedAt: time.Now(),
	}
	mu.Unlock()

	log.Printf("User linked via SSO: state=%s token=%s...", state, token[:8])

	renderSuccess(w, successData{
		TokenPrefix: token[:8],
		LinkedAt:    time.Now().Format("02.01.2006 15:04:05"),
	})
}

// successData — данные для шаблона success_page.html
type successData struct {
	TokenPrefix string
	LinkedAt    string
}

func renderSuccess(w http.ResponseWriter, data successData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "success_page.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// -------- вспомогательные функции --------

// verifyToken проверяет токен у нашего бэкенда (GET /api/v1/sso/me)
func verifyToken(token string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/sso/me?token=%s", cfg.PhoneTokensAPI, url.QueryEscape(token)))
	if err != nil {
		return false, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	valid, _ := result["valid"].(bool)
	return valid, nil
}

// -------- main --------

func main() {
	err := LoadEnv()
	if err != nil {
		log.Fatalf("error loading .env file %v", err)
	}
	cfg = loadConfig()
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/start-sso", startSSOHandler)
	mux.HandleFunc("/callback", callbackHandler)

	log.Printf("External service demo listening on %s", cfg.ListenAddr)
	log.Printf("  Agent ID:          %s", cfg.AgentID)
	log.Printf("  Phone Tokens API:  %s", cfg.PhoneTokensAPI)
	log.Printf("  Callback base URL: %s", cfg.CallbackBaseURL)
	log.Printf("")
	log.Printf("Open http://localhost%s in your browser", cfg.ListenAddr)

	if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
		log.Fatal(err)
	}
}
