package calls

import (
	"context"
	"errors"
	"net/url"
	"testing"
)

func TestConnectClientWithUserViaProxy(t *testing.T) {
	provider := &stubProvider{
		response: &CallbackResponse{
			Status: "success",
			From:   "79990001122",
			To:     "79990002233",
			Time:   1730000000,
		},
	}

	svc := NewService(provider)

	resp, err := svc.ConnectClientWithUserViaProxy(context.Background(), ConnectInput{
		ClientNumber: "79990001122",
		UserNumber:   "79990002233",
		ProxyNumber:  "100",
		Predicted:    true,
	})
	if err != nil {
		t.Fatalf("ConnectClientWithUserViaProxy returned error: %v", err)
	}

	if resp.Status != "success" {
		t.Fatalf("unexpected status: %s", resp.Status)
	}
	if provider.lastRequest.From != "79990001122" {
		t.Fatalf("unexpected from: %s", provider.lastRequest.From)
	}
	if provider.lastRequest.To != "79990002233" {
		t.Fatalf("unexpected to: %s", provider.lastRequest.To)
	}
	if provider.lastRequest.SIP != "100" {
		t.Fatalf("unexpected proxy number: %s", provider.lastRequest.SIP)
	}
	if !provider.lastRequest.Predicted {
		t.Fatalf("expected predicted request")
	}
}

func TestConnectClientWithUserViaProxyValidation(t *testing.T) {
	svc := NewService(&stubProvider{})

	_, err := svc.ConnectClientWithUserViaProxy(context.Background(), ConnectInput{
		UserNumber:  "79990002233",
		ProxyNumber: "100",
	})
	if !errors.Is(err, ErrClientNumberRequired) {
		t.Fatalf("expected ErrClientNumberRequired, got %v", err)
	}

	_, err = svc.ConnectClientWithUserViaProxy(context.Background(), ConnectInput{
		ClientNumber: "79990001122",
	})
	if !errors.Is(err, ErrUserNumberRequired) {
		t.Fatalf("expected ErrUserNumberRequired, got %v", err)
	}
}

func TestConnectClientWithUserWithoutProxy(t *testing.T) {
	provider := &stubProvider{
		response: &CallbackResponse{
			Status: "success",
		},
	}

	svc := NewService(provider)

	_, err := svc.ConnectClientWithUserViaProxy(context.Background(), ConnectInput{
		ClientNumber: "79990001122",
		UserNumber:   "79990002233",
	})
	if err != nil {
		t.Fatalf("ConnectClientWithUserViaProxy returned error: %v", err)
	}
	if provider.lastRequest.SIP != "" {
		t.Fatalf("expected empty SIP for call without proxy, got %s", provider.lastRequest.SIP)
	}
}

func TestConnectClientWithUserViaProxyWithoutProvider(t *testing.T) {
	svc := NewService(nil)

	_, err := svc.ConnectClientWithUserViaProxy(context.Background(), ConnectInput{
		ClientNumber: "79990001122",
		UserNumber:   "79990002233",
		ProxyNumber:  "100",
	})
	if !errors.Is(err, ErrProviderNotConfigured) {
		t.Fatalf("expected ErrProviderNotConfigured, got %v", err)
	}
}

func TestHandleProviderCallbackDelegates(t *testing.T) {
	provider := &stubCallbackProvider{
		stubProvider: stubProvider{},
		response: &CallbackEvent{
			Event:     "NOTIFY_START",
			PBXCallID: "pbx-1",
		},
	}

	svc := NewService(provider)
	form := url.Values{
		"event":       {"NOTIFY_START"},
		"pbx_call_id": {"pbx-1"},
	}

	event, err := svc.HandleProviderCallback(form, "signature")
	if err != nil {
		t.Fatalf("HandleProviderCallback returned error: %v", err)
	}
	if event.Event != "NOTIFY_START" {
		t.Fatalf("unexpected event: %s", event.Event)
	}
	if provider.signature != "signature" {
		t.Fatalf("unexpected signature: %s", provider.signature)
	}
}

func TestHandleProviderCallbackUnsupported(t *testing.T) {
	svc := NewService(&stubProvider{})

	_, err := svc.HandleProviderCallback(url.Values{}, "signature")
	if !errors.Is(err, ErrCallbackNotSupported) {
		t.Fatalf("expected ErrCallbackNotSupported, got %v", err)
	}
}

type stubProvider struct {
	lastRequest CallbackRequest
	response    *CallbackResponse
	err         error
}

func (s *stubProvider) RequestCallback(_ context.Context, req CallbackRequest) (*CallbackResponse, error) {
	s.lastRequest = req
	if s.err != nil {
		return nil, s.err
	}
	if s.response == nil {
		return &CallbackResponse{Status: "success"}, nil
	}
	return s.response, nil
}

type stubCallbackProvider struct {
	stubProvider
	response  *CallbackEvent
	err       error
	signature string
}

func (s *stubCallbackProvider) ParseWebhook(_ url.Values, signature string) (*CallbackEvent, error) {
	s.signature = signature
	if s.err != nil {
		return nil, s.err
	}
	return s.response, nil
}
