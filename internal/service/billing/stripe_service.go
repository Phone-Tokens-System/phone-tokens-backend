package billing

import (
	"fmt"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

var stripeWebhookSecret = "whsec_..." // твой webhook секрет

func (s *BillingService) CreateCheckoutSession(amount float64, agentID string) (string, error) {
	stripe.Key = s.token
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Пополнение баланса"),
					},
					UnitAmount: stripe.Int64(int64(amount * 100)), // в центах
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:8080/success?agent=" + agentID + "&amount=" +
			fmt.Sprintf("%.2f", amount)),
		CancelURL: stripe.String("http://localhost:8080/cancel"),
	}
	params.AddMetadata("agent_id", agentID)

	sSession, err := session.New(params)
	if err != nil {
		return "", err
	}

	return sSession.URL, nil
}
