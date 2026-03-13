package calls

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

var (
	ErrClientNumberRequired  = errors.New("client_number is required")
	ErrUserNumberRequired    = errors.New("user_number is required")
	ErrProviderNotConfigured = errors.New("calls provider is not configured")
	ErrCallbackNotSupported  = errors.New("calls callback is not supported by provider")
)

type service struct {
	provider Provider
}

func NewService(provider Provider) Service {
	return &service{provider: provider}
}

func (s *service) ConnectClientWithUserViaProxy(ctx context.Context, input ConnectInput) (*CallbackResponse, error) {
	if s.provider == nil {
		return nil, ErrProviderNotConfigured
	}

	clientNumber := strings.TrimSpace(input.ClientNumber)
	if clientNumber == "" {
		return nil, ErrClientNumberRequired
	}

	userNumber := strings.TrimSpace(input.UserNumber)
	if userNumber == "" {
		return nil, ErrUserNumberRequired
	}

	return s.provider.RequestCallback(ctx, CallbackRequest{
		From:      clientNumber,
		To:        userNumber,
		SIP:       strings.TrimSpace(input.ProxyNumber),
		Predicted: input.Predicted,
	})
}

func (s *service) HandleProviderCallback(form url.Values, signature string) (*CallbackEvent, error) {
	if s.provider == nil {
		return nil, ErrProviderNotConfigured
	}

	parser, ok := s.provider.(CallbackParser)
	if !ok {
		return nil, ErrCallbackNotSupported
	}

	return parser.ParseWebhook(form, signature)
}
