package sms_service

import (
	"phone-tokens/internal/model"
)

type SmsAdapter interface {
	SendSms(number int, text string) (model.SmsResponse, error)
	GetSmsStatus(id int) (model.SmsStatus, error)
	GetSmsList() ([]model.SmsResponse, error)
}
