package model

/*
*
request to send sms to one user.
*/
type SmsRequest struct {
	Certificate  []byte `json:"certificate"`
	ClientNumber int    `json:"client_number"`
	Text         string `json:"text"`
}
