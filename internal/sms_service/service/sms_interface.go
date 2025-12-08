package service

import (
	"phone-tokens/internal/sms_service/model/interface"
)

type SmsAdapter interface {
	SendSms(number int, text string) (_interface.SmsResponse, error)
	GetSmsStatus(id int) (_interface.SmsStatus, error)
	GetSmsList() ([]map[string]string, error)
}
