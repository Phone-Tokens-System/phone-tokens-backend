package calls

import (
	"context"
	"net/url"
)

type ConnectInput struct {
	ClientNumber string
	UserNumber   string
	ProxyNumber  string
	Predicted    bool
}

type CallbackRequest struct {
	From      string
	To        string
	SIP       string
	Predicted bool
}

type CallbackResponse struct {
	Status string `json:"status"`
	From   string `json:"from"`
	To     string `json:"to"`
	Time   int64  `json:"time"`
}

type CallbackEvent struct {
	Event         string            `json:"event"`
	CallStart     string            `json:"call_start,omitempty"`
	PBXCallID     string            `json:"pbx_call_id,omitempty"`
	CallerID      string            `json:"caller_id,omitempty"`
	CalledDID     string            `json:"called_did,omitempty"`
	Destination   string            `json:"destination,omitempty"`
	Internal      string            `json:"internal,omitempty"`
	Duration      string            `json:"duration,omitempty"`
	Disposition   string            `json:"disposition,omitempty"`
	StatusCode    string            `json:"status_code,omitempty"`
	IsRecorded    string            `json:"is_recorded,omitempty"`
	CallIDWithRec string            `json:"call_id_with_rec,omitempty"`
	Raw           map[string]string `json:"raw,omitempty"`
}

type Provider interface {
	RequestCallback(ctx context.Context, req CallbackRequest) (*CallbackResponse, error)
}

type CallbackParser interface {
	ParseWebhook(form url.Values, signature string) (*CallbackEvent, error)
}

type Service interface {
	ConnectClientWithUserViaProxy(ctx context.Context, input ConnectInput) (*CallbackResponse, error)
	HandleProviderCallback(form url.Values, signature string) (*CallbackEvent, error)
}
