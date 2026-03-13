package novofon

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"phone-tokens/internal/service/calls"
)

const (
	defaultBaseURL   = "https://api.novofon.com"
	defaultTimeout   = 10 * time.Second
	callbackEndpoint = "/v1/request/callback/"
)

var ErrAPI = errors.New("novofon api error")

type Config struct {
	APIKey    string
	APISecret string
	BaseURL   string
	Timeout   time.Duration
}

type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client
}

func NewClient(cfg Config) (*Client, error) {
	apiKey := strings.TrimSpace(cfg.APIKey)
	if apiKey == "" {
		return nil, errors.New("novofon api key is required")
	}

	apiSecret := strings.TrimSpace(cfg.APISecret)
	if apiSecret == "" {
		return nil, errors.New("novofon api secret is required")
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (c *Client) RequestCallback(ctx context.Context, req calls.CallbackRequest) (*calls.CallbackResponse, error) {
	from := strings.TrimSpace(req.From)
	if from == "" {
		return nil, errors.New("from is required")
	}

	to := strings.TrimSpace(req.To)
	if to == "" {
		return nil, errors.New("to is required")
	}

	params := url.Values{}
	params.Set("from", from)
	params.Set("to", to)
	params.Set("format", "json")

	sip := strings.TrimSpace(req.SIP)
	if sip != "" {
		params.Set("sip", sip)
	}

	if req.Predicted {
		params.Set("predicted", "1")
	}

	query := params.Encode()
	signature := c.makeSignature(callbackEndpoint, query)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseURL+callbackEndpoint+"?"+query,
		nil,
	)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", c.apiKey+":"+signature)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("%w: http %d: %s", ErrAPI, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data := make(map[string]any)
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode novofon response: %w", err)
	}

	status := toString(data["status"])
	if strings.EqualFold(status, "error") {
		message := toString(data["message"])
		if message == "" {
			message = "unknown error"
		}
		return nil, fmt.Errorf("%w: %s", ErrAPI, message)
	}

	return &calls.CallbackResponse{
		Status: status,
		From:   toString(data["from"]),
		To:     toString(data["to"]),
		Time:   toInt64(data["time"]),
	}, nil
}

func (c *Client) makeSignature(method, paramsStr string) string {
	// Novofon expects base64-encoded hex hmac sha1 digest, same as official SDK.
	stringToSign := method + paramsStr + md5Hex(paramsStr)
	return encodeSignature(stringToSign, c.apiSecret)
}

func md5Hex(value string) string {
	sum := md5.Sum([]byte(value))
	return hex.EncodeToString(sum[:])
}

func encodeSignature(payload, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	sha1Hex := hex.EncodeToString(mac.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha1Hex))
}

func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.Itoa(v)
	default:
		if v == nil {
			return ""
		}
		return fmt.Sprintf("%v", v)
	}
}

func toInt64(value any) int64 {
	switch v := value.(type) {
	case json.Number:
		i, err := v.Int64()
		if err == nil {
			return i
		}
		f, err := v.Float64()
		if err == nil {
			return int64(f)
		}
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i
		}
	}
	return 0
}
