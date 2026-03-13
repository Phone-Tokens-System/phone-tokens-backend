package novofon

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

func TestParseWebhookNotifyStart(t *testing.T) {
	client, err := NewClient(Config{
		APIKey:    "key",
		APISecret: "secret",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	form := url.Values{
		"event":       {EventNotifyStart},
		"call_start":  {"1730000000"},
		"pbx_call_id": {"pbx-1"},
		"caller_id":   {"79990001122"},
		"called_did":  {"79990002233"},
	}
	signature := expectedWebhookSignature("79990001122799900022331730000000", "secret")

	event, err := client.ParseWebhook(form, signature)
	if err != nil {
		t.Fatalf("ParseWebhook returned error: %v", err)
	}
	if event.Event != EventNotifyStart {
		t.Fatalf("unexpected event: %s", event.Event)
	}
	if event.PBXCallID != "pbx-1" {
		t.Fatalf("unexpected pbx_call_id: %s", event.PBXCallID)
	}
	if event.CallerID != "79990001122" {
		t.Fatalf("unexpected caller_id: %s", event.CallerID)
	}
}

func TestParseWebhookMissingSignature(t *testing.T) {
	client, err := NewClient(Config{
		APIKey:    "key",
		APISecret: "secret",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	_, err = client.ParseWebhook(url.Values{"event": {EventNotifyStart}}, "")
	if !errors.Is(err, ErrMissingWebhookSignature) {
		t.Fatalf("expected ErrMissingWebhookSignature, got %v", err)
	}
}

func TestParseWebhookInvalidSignature(t *testing.T) {
	client, err := NewClient(Config{
		APIKey:    "key",
		APISecret: "secret",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	form := url.Values{
		"event":       {EventNotifyAnswer},
		"caller_id":   {"79990001122"},
		"destination": {"100"},
		"call_start":  {"1730000000"},
	}

	_, err = client.ParseWebhook(form, "bad-signature")
	if !errors.Is(err, ErrInvalidWebhookSignature) {
		t.Fatalf("expected ErrInvalidWebhookSignature, got %v", err)
	}
}

func TestParseWebhookUnsupportedEvent(t *testing.T) {
	client, err := NewClient(Config{
		APIKey:    "key",
		APISecret: "secret",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	form := url.Values{
		"event": {"UNKNOWN"},
	}
	signature := expectedWebhookSignature("whatever", "secret")

	_, err = client.ParseWebhook(form, signature)
	if !errors.Is(err, ErrUnsupportedWebhookEvent) {
		t.Fatalf("expected ErrUnsupportedWebhookEvent, got %v", err)
	}
}

func expectedWebhookSignature(payload, secret string) string {
	return encodeSignature(payload, secret)
}
