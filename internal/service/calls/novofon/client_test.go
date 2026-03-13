package novofon

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"phone-tokens/internal/service/calls"
)

func TestRequestCallbackSuccess(t *testing.T) {
	const (
		expectedQuery = "format=json&from=79990001122&predicted=1&sip=100&to=79990002233"
		apiKey        = "user-key"
		apiSecret     = "secret-key"
	)

	var gotAuthHeader string
	var gotPath string
	var gotQuery string

	client, err := NewClient(Config{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   "https://api.novofon.com",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("expected GET method, got %s", r.Method)
			}
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			gotAuthHeader = r.Header.Get("Authorization")
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","from":79990001122,"to":"79990002233","time":1730000000}`)),
				Request:    r,
			}, nil
		}),
	}

	resp, err := client.RequestCallback(context.Background(), calls.CallbackRequest{
		From:      "79990001122",
		To:        "79990002233",
		SIP:       "100",
		Predicted: true,
	})
	if err != nil {
		t.Fatalf("RequestCallback returned error: %v", err)
	}

	expectedAuthorization := apiKey + ":" + expectedSignature(callbackEndpoint, expectedQuery, apiSecret)
	if gotAuthHeader != expectedAuthorization {
		t.Fatalf("unexpected Authorization header: %s", gotAuthHeader)
	}
	if gotPath != callbackEndpoint {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotQuery != expectedQuery {
		t.Fatalf("unexpected query: %s", gotQuery)
	}

	if resp.Status != "success" {
		t.Fatalf("unexpected status: %s", resp.Status)
	}
	if resp.From != "79990001122" {
		t.Fatalf("unexpected from: %s", resp.From)
	}
	if resp.To != "79990002233" {
		t.Fatalf("unexpected to: %s", resp.To)
	}
	if resp.Time != 1730000000 {
		t.Fatalf("unexpected time: %d", resp.Time)
	}
}

func TestRequestCallbackReturnsAPIError(t *testing.T) {
	client, err := NewClient(Config{
		APIKey:    "user-key",
		APISecret: "secret-key",
		BaseURL:   "https://api.novofon.com",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"status":"error","message":"wrong number"}`)),
				Request:    r,
			}, nil
		}),
	}

	_, err = client.RequestCallback(context.Background(), calls.CallbackRequest{
		From: "79990001122",
		To:   "bad-number",
		SIP:  "100",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrAPI) {
		t.Fatalf("expected ErrAPI, got %v", err)
	}
}

func TestRequestCallbackHTTPStatusError(t *testing.T) {
	client, err := NewClient(Config{
		APIKey:    "user-key",
		APISecret: "secret-key",
		BaseURL:   "https://api.novofon.com",
		Timeout:   time.Second,
	})
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Body:       io.NopCloser(strings.NewReader(`bad gateway`)),
				Request:    r,
			}, nil
		}),
	}

	_, err = client.RequestCallback(context.Background(), calls.CallbackRequest{
		From: "79990001122",
		To:   "79990002233",
		SIP:  "100",
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, ErrAPI) {
		t.Fatalf("expected ErrAPI, got %v", err)
	}
}

func expectedSignature(method, params, secret string) string {
	md5Sum := md5.Sum([]byte(params))
	stringToSign := method + params + hex.EncodeToString(md5Sum[:])
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(stringToSign))
	sha1Hex := hex.EncodeToString(mac.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha1Hex))
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme == "" {
		req.URL.Scheme = "https"
	}
	if req.URL.Host == "" {
		req.URL.Host = "api.novofon.com"
	}
	return f(req)
}
