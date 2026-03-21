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
	AgentID     string `json:"agent_id,omitempty" gorm:"-"`
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
	Id           string    `json:"id,omitempty" gorm:"column:id;not null; primaryKey;autoIncrement"`
	ExternalId   string    `json:"external_id" gorm:"column:external_id;not null"`
	From         string    `json:"from,omitempty" gorm:"column:from_number"`
	Token        string    `json:"token,omitempty" gorm:"column:token;not null"`
	Text         string    `json:"text" gorm:"not null"`
	Status       int       `json:"status" gorm:"not null"`
	ExtendStatus string    `json:"extend_status,omitempty" gorm:"column:extended_status"`
	Cost         float64   `json:"cost"`
	DateCreated  int       `json:"date_created,omitempty"`
	DateSent     int       `json:"date_sent,omitempty"`
	Raw          []byte    `json:"raw,omitempty" gorm:"type:text"`
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
