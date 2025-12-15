package model

import "github.com/google/uuid"

/*
*
request to send sms to one user.
*/
type SmsRequest struct {
	ServiceName string `json:"service_name"`
	Certificate string `json:"certificate"`
	ClientToken string `json:"client_token"`
	Text        string `json:"text"`
}

/*
*
struct for mass-request for sms.
a lot of clients - one text
*/
type GroupSMSRequest struct {
	ServiceName   string `json:"service_name"`
	Certificate   []byte `json:"certificate"`
	ClientNumbers []int  `json:"client_numbers"`
	Text          string `json:"text"`
}

type SmsResponse struct {
	ServiceName  string    `json:"service_name" gorm:"not null"`
	ServiceId    uuid.UUID `json:"service_id" gorm:"not null"`
	Id           string    `json:"id,omitempty" gorm:"column:external_id;not null"`
	From         string    `json:"from,omitempty" gorm:"column:from_number"`
	Number       string    `json:"number" gorm:"not null"`
	Text         string    `json:"text" gorm:"not null"`
	Status       int       `json:"status" gorm:"not null"`
	ExtendStatus string    `json:"extend_status,omitempty"`
	Cost         float64   `json:"cost"`
	DateCreated  int       `json:"date_created,omitempty"`
	DateSent     int       `json:"date_sent,omitempty"`
	Raw          any       `json:"raw,omitempty" gorm:"raw;not null"`
}

type SmsStatus struct {
	ServiceName  string `json:"service_name"`
	Number       string `json:"number"`
	Status       int    `json:"status,omitempty"`
	ExtendStatus string `json:"extend_status,omitempty"`
}

func (SmsResponse) TableName() string {
	return "sms"
}
