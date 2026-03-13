package novofon

import (
	"errors"
	"net/url"
	"strings"

	"phone-tokens/internal/service/calls"
)

const (
	EventNotifyStart    = "NOTIFY_START"
	EventNotifyInternal = "NOTIFY_INTERNAL"
	EventNotifyAnswer   = "NOTIFY_ANSWER"
	EventNotifyEnd      = "NOTIFY_END"
	EventNotifyOutStart = "NOTIFY_OUT_START"
	EventNotifyOutEnd   = "NOTIFY_OUT_END"
	EventNotifyRecord   = "NOTIFY_RECORD"
	EventNotifyIVR      = "NOTIFY_IVR"
)

var (
	ErrMissingWebhookSignature = errors.New("missing callback signature")
	ErrInvalidWebhookSignature = errors.New("invalid callback signature")
	ErrUnsupportedWebhookEvent = errors.New("unsupported callback event")
)

func (c *Client) ParseWebhook(form url.Values, signature string) (*calls.CallbackEvent, error) {
	signature = strings.TrimSpace(signature)
	if signature == "" {
		return nil, ErrMissingWebhookSignature
	}

	eventType := strings.TrimSpace(form.Get("event"))
	if eventType == "" {
		return nil, ErrUnsupportedWebhookEvent
	}

	signaturePayload, err := webhookSignaturePayload(eventType, form)
	if err != nil {
		return nil, err
	}

	expected := encodeSignature(signaturePayload, c.apiSecret)
	if signature != expected {
		return nil, ErrInvalidWebhookSignature
	}

	return &calls.CallbackEvent{
		Event:         eventType,
		CallStart:     form.Get("call_start"),
		PBXCallID:     form.Get("pbx_call_id"),
		CallerID:      form.Get("caller_id"),
		CalledDID:     form.Get("called_did"),
		Destination:   form.Get("destination"),
		Internal:      form.Get("internal"),
		Duration:      form.Get("duration"),
		Disposition:   form.Get("disposition"),
		StatusCode:    form.Get("status_code"),
		IsRecorded:    form.Get("is_recorded"),
		CallIDWithRec: form.Get("call_id_with_rec"),
		Raw:           firstValues(form),
	}, nil
}

func webhookSignaturePayload(eventType string, form url.Values) (string, error) {
	switch eventType {
	case EventNotifyStart, EventNotifyInternal, EventNotifyEnd, EventNotifyIVR:
		return form.Get("caller_id") + form.Get("called_did") + form.Get("call_start"), nil
	case EventNotifyAnswer:
		return form.Get("caller_id") + form.Get("destination") + form.Get("call_start"), nil
	case EventNotifyOutStart, EventNotifyOutEnd:
		return form.Get("internal") + form.Get("destination") + form.Get("call_start"), nil
	case EventNotifyRecord:
		return form.Get("pbx_call_id") + form.Get("call_id_with_rec"), nil
	default:
		return "", ErrUnsupportedWebhookEvent
	}
}

func firstValues(values url.Values) map[string]string {
	result := make(map[string]string, len(values))
	for key, val := range values {
		if len(val) == 0 {
			result[key] = ""
			continue
		}
		result[key] = val[0]
	}
	return result
}
