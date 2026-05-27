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
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// -------- конфигурация --------

var (
	agentID         = getenv("AGENT_ID", "00000000-0000-0000-0000-000000000001")
	phoneTokensAPI  = getenv("PHONE_TOKENS_API", "http://localhost:8080")
	listenAddr      = getenv("LISTEN_ADDR", ":9090")
	callbackBaseURL = getenv("CALLBACK_BASE_URL", "http://localhost:9090")
)

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
	html := `<!DOCTYPE html>
<html lang="ru">
<head><meta charset="utf-8"><title>Мой сервис</title>
<style>
  body { font-family: sans-serif; max-width: 600px; margin: 60px auto; }
  .btn { display: inline-block; padding: 12px 24px; background: #6c47ff;
         color: white; border-radius: 8px; text-decoration: none; font-size: 16px; }
  .btn:hover { background: #5534d4; }
  pre { background: #f4f4f4; padding: 12px; border-radius: 6px; overflow-x: auto; }
</style>
</head>
<body>
  <h1>Добро пожаловать в Мой Сервис</h1>
  <p>Для регистрации используйте Phone Tokens — вам не нужно вводить номер телефона здесь:</p>
  <a class="btn" href="/start-sso">Зарегистрироваться через Phone Tokens</a>
  <hr>
  <p><small>После привязки вы сможете получать SMS от нас без раскрытия реального номера.</small></p>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
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
	redirectURI := callbackBaseURL + "/callback"
	ssoURL := fmt.Sprintf(
		"%s/api/v1/sso/authorize?agent_id=%s&redirect_uri=%s&state=%s",
		phoneTokensAPI,
		url.QueryEscape(agentID),
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

	// Показываем страницу успеха
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head><meta charset="utf-8"><title>Успешная привязка</title>
<style>
  body { font-family: sans-serif; max-width: 600px; margin: 60px auto; }
  .card { background: #f0fff0; border: 1px solid #6c6; border-radius: 8px; padding: 20px; }
  code { background: #eee; padding: 2px 6px; border-radius: 4px; }
  .btn { display: inline-block; margin-top: 20px; padding: 10px 20px;
         background: #6c47ff; color: white; border-radius: 6px;
         text-decoration: none; }
</style>
</head>
<body>
  <h1>✅ Привязка выполнена!</h1>
  <div class="card">
    <p>Ваш телефонный токен успешно привязан к нашему сервису.</p>
    <p>Токен (первые 8 символов): <code>%s...</code></p>
    <p>Привязан: <code>%s</code></p>
    <p>Теперь мы можем отправлять вам SMS через Phone Tokens систему,
       не зная вашего реального номера.</p>
  </div>
  <a class="btn" href="/send-test-sms?state=%s">Отправить тестовое SMS</a>
  <a class="btn" style="background:#888; margin-left:10px;" href="/">На главную</a>
</body>
</html>`, token[:8], time.Now().Format("02.01.2006 15:04:05"), state)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// GET /send-test-sms?state=<state> — демонстрационная отправка SMS через наш API
func sendTestSMSHandler(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")

	mu.Lock()
	user, ok := users[state]
	mu.Unlock()

	if !ok || user == nil {
		http.Error(w, "user not found or not linked", http.StatusBadRequest)
		return
	}

	// Отправляем SMS через API нашей системы
	// NOTE: для реального запроса нужен сертификат агента (mTLS).
	// Здесь для простоты используем Bearer-токен агента.
	// В продакшене — client certificate через https.
	result, err := sendSMSViaAPI(user.Token, "Привет! Это тестовое SMS от Моего Сервиса через Phone Tokens 🎉")
	if err != nil {
		http.Error(w, "SMS send failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="ru"><head><meta charset="utf-8"><title>SMS отправлено</title></head>
<body>
  <h1>📱 SMS отправлено!</h1>
  <p>Результат: <pre>%s</pre></p>
  <a href="/">← Назад</a>
</body></html>`, result)
}

// -------- вспомогательные функции --------

// verifyToken проверяет токен у нашего бэкенда (GET /api/v1/sso/me)
func verifyToken(token string) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/sso/me?token=%s", phoneTokensAPI, url.QueryEscape(token)))
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

// sendSMSViaAPI отправляет SMS через POST /api/v1/sms/send
// В реальном сервисе здесь будет mTLS с сертификатом агента.
// Для демо используется агентский JWT (получить через POST /api/v1/login).
func sendSMSViaAPI(clientToken, text string) (string, error) {
	agentJWT := os.Getenv("AGENT_JWT") // JWT агента, полученный через /api/v1/login
	if agentJWT == "" {
		return "", fmt.Errorf("AGENT_JWT env var not set — cannot call SMS API without auth")
	}

	payload := map[string]string{
		"client_token": clientToken,
		"text":         text,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", phoneTokensAPI+"/api/v1/sms/send", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+agentJWT)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}
	return string(respBody), nil
}

// -------- main --------

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/start-sso", startSSOHandler)
	mux.HandleFunc("/callback", callbackHandler)
	mux.HandleFunc("/send-test-sms", sendTestSMSHandler)

	log.Printf("External service demo listening on %s", listenAddr)
	log.Printf("  Agent ID:          %s", agentID)
	log.Printf("  Phone Tokens API:  %s", phoneTokensAPI)
	log.Printf("  Callback base URL: %s", callbackBaseURL)
	log.Printf("")
	log.Printf("Open http://localhost%s in your browser", listenAddr)

	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatal(err)
	}
}
