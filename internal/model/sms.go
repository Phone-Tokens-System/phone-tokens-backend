package model

/*
*
request to send sms_service to one user.
*/
type SmsRequest struct {
	Certificate  []byte `json:"certificate"`
	ClientNumber int    `json:"client_number"`
	Text         string `json:"text"`
}

//	type HlrCheck struct {
//	   Id              int
//	   Number          string
//	   HlrStatus       int
//	   ExtendHlrStatus string
//	}
type SmsResponse struct {
	Id           int     `json:"id,omitempty"`
	From         string  `json:"from,omitempty"`
	Number       string  `json:"number"`
	Text         string  `json:"text"`
	Status       int     `json:"status"`
	ExtendStatus string  `json:"extend_status,omitempty"`
	Cost         float64 `json:"cost"`
	DateCreated  int     `json:"date_created,omitempty"`
	DateSent     int     `json:"date_sent,omitempty"`
	Raw          any     `json:"raw,omitempty"`
}

//	type SendSms struct {
//	   Id           int
//	   From         string
//	   Number       string
//	   Text         string
//	   Status       int
//	   ExtendStatus string
//	   Channel      string
//	   Cost         float64
//	   DateCreate   int
//	   DateSend     int
//	}
type SmsStatus struct {
	Number       string `json:"number"`
	Status       int    `json:"status,omitempty"`
	ExtendStatus string `json:"extend_status,omitempty"`
}

//type SmsInfo struct {
//	Id           int    `json:"id,omitempty"`
//	Number       string `json:"number"`
//	Status       int    `json:"status,omitempty"`
//	ExtendStatus string `json:"extend_status,omitempty"`
//	TimeCreated  int    `json:"time_created,omitempty"`
//	TimeSent     int    `json:"time_sent,omitempty"`
//}
