package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"phone-tokens/internal/service/billing"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
)

type BillingHandler struct {
	BillingService *billing.BillingService
	WebhookSecret  string
}

var stripeWebhookSecret = "whsec_..." // твой webhook секрет

func NewBillingHandler(billingService *billing.BillingService) *BillingHandler {
	return &BillingHandler{BillingService: billingService}
}

// POST /create-checkout
// Body: { "agent_id": "42", "amount": 10.0 }
// Response: { "checkout_url": "https://checkout.stripe.com/..." }
// → window.location.href = checkout_url

// TopBalance godoc
// @Summary Создать Stripe Checkout session для пополнения баланса
// @Description Создает checkout-сессию на Stripe и возвращает URL для редиректа
// @Security BearerAuth
// @Tags Billing
// @Accept json
// @Produce json
// @Param request body dto.TopUpRequest true "Request payload"
// @Success 200 {object} map[string]string "URL для редиректа на Stripe Checkout"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /api/v1/billing/balance [post]
func (h *BillingHandler) TopBalance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AgentID string  `json:"agent_id"`
		Amount  float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.BillingService.CreateCheckoutSession(req.Amount, req.AgentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// возвращаем JSON с URL для редиректа на Stripe
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"checkout_url": url,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

// StripeWebhookHandler godoc
// @Summary Обработка Stripe Webhook
// @Description Получает уведомления о событиях Stripe, например успешной оплате checkout session
// @Tags Billing
// @Accept json
// @Produce json
// @Param stripe-signature header string true "Stripe Signature"
// @Success 200 "Webhook обработан успешно"
// @Failure 400 "Некорректный payload или подпись"
// @Router /api/v1/billing/webhook [post]
func (h *BillingHandler) StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	payload, _ := io.ReadAll(r.Body)
	sigHeader := r.Header.Get("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, sigHeader, stripeWebhookSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		agentID := session.Metadata["agent_id"]
		amount := float64(session.AmountTotal) / 100.0
		_ = h.BillingService.TopUpBalance(r.Context(), agentID, amount)
	}

	w.WriteHeader(http.StatusOK)
}

// GetBalanceHandler godoc
// @Summary Получить текущий баланс агента
// @Description Возвращает float64 баланс агента по его ID
// @Security BearerAuth
// @Tags Billing
// @Accept json
// @Produce json
// @Param agent_id query string true "ID агента"
// @Success 200 {object} map[string]float64 "Баланс агента"
// @Failure 404 {object} map[string]string "Агент не найден"
// @Router /api/v1/billing/balance [get]
func (h *BillingHandler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	agentID := r.URL.Query().Get("agent_id")
	balance, err := h.BillingService.GetBalance(r.Context(), agentID) // метод возвращает float64
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
	if err != nil {
		fmt.Println(err)
		return
	}
}
