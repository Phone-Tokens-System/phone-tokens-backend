package model

import "github.com/google/uuid"

/*
*
request to send sms_service to one user.
*/
type SmsRequest struct {
	ServiceName string `json:"service_name"`
	Certificate string `json:"certificate"`
	ClientToken string `json:"client_token"`
	Text        string `json:"text"`
}

/*
*
struct for mass-request for sms_service.
a lot of clients - one text
*/
type GroupSMSRequest struct {
	ServiceName   string `json:"service_name"`
	Certificate   []byte `json:"certificate"`
	ClientNumbers []int  `json:"client_numbers"`
	Text          string `json:"text"`
}

type SmsResponse struct {
	ServiceName  string    `json:"service_name"`
	ServiceId    uuid.UUID `json:"service_id"`
	Id           int       `json:"id,omitempty"`
	From         string    `json:"from,omitempty"`
	Number       string    `json:"number"`
	Text         string    `json:"text"`
	Status       int       `json:"status"`
	ExtendStatus string    `json:"extend_status,omitempty"`
	Cost         float64   `json:"cost"`
	DateCreated  int       `json:"date_created,omitempty"`
	DateSent     int       `json:"date_sent,omitempty"`
	Raw          any       `json:"raw,omitempty"`
}

type SmsStatus struct {
	ServiceName  string `json:"service_name"`
	Number       string `json:"number"`
	Status       int    `json:"status,omitempty"`
	ExtendStatus string `json:"extend_status,omitempty"`
}
