package dto

import (
	"phone-tokens/internal/model"

	"github.com/google/uuid"
)

type SmsLog struct {
	ServiceName string    `json:"service_name"`
	ServiceId   uuid.UUID `json:"service_id"`
	From        string    `json:"from"`
	Token       string    `json:"token"`
	Text        string    `json:"text"`
	Status      string    `json:"status"`
	Cost        float64   `json:"cost"`
	DateCreated int       `json:"date_created"`
	DateSent    int       `json:"date_sent"`
}

func ToSmsLog(sms model.SmsResponse) *SmsLog {
	return &SmsLog{
		ServiceName: sms.ServiceName,
		ServiceId:   sms.ServiceId,
		From:        sms.From,
		Token:       sms.Token,
		Text:        sms.Text,
		Cost:        sms.Cost,
		DateCreated: sms.DateCreated,
		DateSent:    sms.DateSent,
	}
}
