package sms

import (
	"context"
	"phone-tokens/internal/model"
)

// TODO: service
type Repository interface {
	SaveSms(ctc context.Context, smsResponse model.SmsResponse) error
	GetAllSms(ctx context.Context) ([]model.SmsResponse, error)
	GetSmsByServiceId(ctx context.Context, serviceId string) ([]model.SmsResponse, error)
	GetSmsByToken(ctx context.Context, token string) ([]model.SmsResponse, error)
	GetSmsByServiceName(ctx context.Context, serviceName string) ([]model.SmsResponse, error)
}
