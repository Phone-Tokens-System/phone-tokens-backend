package service

import "phone-tokens/internal/sms_service/model"

type SmsService interface {
	SendSms(number int, text string) (model.SmsResponse, error)
}
