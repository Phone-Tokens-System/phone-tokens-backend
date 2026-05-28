package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"phone-tokens/internal/adapter/dto"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/billing"
	"strings"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
)

type BillingHandler struct {
	BillingService *billing.BillingService
	WebhookSecret  string
}

func NewBillingHandler(billingService *billing.BillingService, secret string) *BillingHandler {
	return &BillingHandler{BillingService: billingService, WebhookSecret: secret}
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
// @Param request body dto.BalanceRequest true "Request payload"
// @Success 200 {object} map[string]string "URL для редиректа на Stripe Checkout"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /api/v1/billing/balance [post]
func (h *BillingHandler) TopBalance(w http.ResponseWriter, r *http.Request) {
	var req dto.BalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	url, err := h.BillingService.CreateCheckoutSession(req.Amount, req.AgentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"checkout_url": url,
	})
	// возвращаем JSON с URL для редиректа на Stripe
	//w.Header().Set("Content-Type", "application/json")
	//err = json.NewEncoder(w).Encode(map[string]string{
	//    "checkout_url": url,
	//})
	//if err != nil {
	//    fmt.Println(err)
	//    return
	//}
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

	event, err := webhook.ConstructEventWithOptions(
		payload,
		sigHeader,
		h.WebhookSecret,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("stripe webhook bad request error", "error", err)
		return
	}

	if event.Type == "checkout.session.completed" {

		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			slog.Error("stripe webhook bad request error", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		agentID := session.Metadata["agent_id"]
		amount := float64(session.AmountTotal) / 100.0
		err = h.BillingService.TopUpBalance(r.Context(), agentID, amount)
		if err != nil {
			slog.Error("stripe webhook top up balance error", "error", err)
		}
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
// @Param agent_id path string true "ID агента"
// @Success 200 {object} map[string]float64 "Баланс агента"
// @Failure 404 {object} map[string]string "Агент не найден"
// @Router /api/v1/billing/balance [get]
func (h *BillingHandler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	agentID := strings.TrimSpace(r.PathValue("agent_id"))
	if agentID == "" {
		agentID = strings.TrimSpace(r.URL.Query().Get("agent_id"))
	}
	if agentID == "" {
		http.Error(w, "agent_id is required", http.StatusBadRequest)
		return
	}

	balance, err := h.BillingService.GetBalance(r.Context(), agentID) // метод возвращает float64
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"balance": balance})
}

// SeePackageOptionsHandler godoc
// @Summary Получить список доступных пакетов
// @Description Возвращает список всех пакетов, которые агент может купить
// @Tags Billing
// @Accept json
// @Produce json
// @Success 200 {array} model.Package
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/packages [get]
func (h *BillingHandler) SeePackageOptionsHandler(w http.ResponseWriter, r *http.Request) {
	pkgs, err := h.BillingService.GetPackages(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, pkgs)
}

// BuyPackageHandler godoc
// @Summary Купить пакет
// @Description Агент покупает пакет услуг (например SMS или звонки)
// @Tags Billing
// @Accept json
// @Produce json
// @Param agent_id path string true "ID агента"
// @Param request body dto.BuyPackageRequest true "Данные покупки пакета"
// @Success 200 {object} map[string]string "package purchased"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/agents/{agent_id}/packages [post]
func (h *BillingHandler) BuyPackageHandler(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("agent_id")

	var req dto.BuyPackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := h.BillingService.AddAgentPkg(r.Context(), req.PkgID, agentID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, billing.ErrNotEnoughBalance) {
			status = http.StatusPaymentRequired
		} else if err.Error() == "package not found" {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetTransactionsHandler godoc
// @Summary Получить историю транзакций агента
// @Description Возвращает список всех транзакций (пополнений и списаний) для агента
// @Tags Billing
// @Accept json
// @Produce json
// @Param agent_id path string true "ID агента"
// @Success 200 {array} model.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/agents/{agent_id}/transactions [get]
func (h *BillingHandler) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("agent_id")
	if agentID == "" {
		http.Error(w, "agent_id is required", http.StatusBadRequest)
		return
	}
	txns, err := h.BillingService.GetTransactions(r.Context(), agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, txns)
}

// CreatePackageHandler godoc
// @Summary Создать тарифный пакет (только для администратора)
// @Description Создаёт новый пакет (например 100 SMS в месяц за 500 руб)
// @Tags Billing
// @Accept json
// @Produce json
// @Param request body dto.CreatePackageRequest true "Данные пакета"
// @Success 201 {object} model.Package
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/packages [post]
func (h *BillingHandler) CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Units <= 0 || req.Price <= 0 {
		http.Error(w, "name, units and price are required", http.StatusBadRequest)
		return
	}

	pkg := &model.Package{
		Name:         req.Name,
		Service:      model.ServiceType(req.Service),
		Units:        req.Units,
		Price:        req.Price,
		DurationDays: req.DurationDays,
		Description:  req.Description,
	}

	if err := h.BillingService.CreatePackage(r.Context(), pkg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, pkg)
}

// DeletePackageHandler godoc
// @Summary Удалить тарифный пакет (только для администратора)
// @Description Удаляет пакет по ID
// @Tags Billing
// @Param pkg_id path string true "ID пакета"
// @Success 204 "Пакет удалён"
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/packages/{pkg_id} [delete]
func (h *BillingHandler) DeletePackageHandler(w http.ResponseWriter, r *http.Request) {
	pkgID := r.PathValue("pkg_id")
	if pkgID == "" {
		http.Error(w, "pkg_id is required", http.StatusBadRequest)
		return
	}
	if err := h.BillingService.DeletePackage(r.Context(), pkgID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SeeAgentPackagesHandler godoc
// @Summary Получить пакеты агента
// @Description Возвращает все пакеты, купленные агентом
// @Tags Billing
// @Accept json
// @Produce json
// @Param agent_id path string true "ID агента"
// @Success 200 {array} model.AgentPackages
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/agents/{agent_id}/packages [get]
func (h *BillingHandler) SeeAgentPackagesHandler(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("agent_id")
	pkgs, err := h.BillingService.GetPackagesByAgentId(r.Context(), agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, pkgs)
}
