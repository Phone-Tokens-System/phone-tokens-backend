package service

import "phone-tokens/internal/sms_service/model"

type SmsAdapter interface {
	SendSms(number int, text string) (model.SmsResponse, error)
}
