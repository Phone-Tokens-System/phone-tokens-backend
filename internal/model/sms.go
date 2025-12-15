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
	ServiceName  string    `json:"service_name" gorm:"column:service_name;not null"`
	ServiceId    uuid.UUID `json:"service_id" gorm:"column:service_id;not null"`
	Id           string    `json:"id,omitempty" gorm:"column:external_id;not null"`
	From         string    `json:"from,omitempty" gorm:"column:from_number"`
	Number       string    `json:"number" gorm:"column:number;not null"`
	Text         string    `json:"text" gorm:"column:text;not null"`
	Status       int       `json:"status" gorm:"column:status;not null"`
	ExtendStatus string    `json:"extend_status,omitempty" gorm:"column:extend_status"`
	Cost         float64   `json:"cost" gorm:"column:cost"`
	DateCreated  int       `json:"date_created,omitempty" gorm:"column:date_created"`
	DateSent     int       `json:"date_sent,omitempty" gorm:"column:date_sent"`
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
