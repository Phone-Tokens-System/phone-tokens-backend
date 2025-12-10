package model

/*
*
struct for mass-request for sms_service.
a lot of clients - one text
*/
type GroupSMSRequest struct {
	Certificate   []byte `json:"certificate"`
	ClientNumbers []int  `json:"client_numbers"`
	Text          string `json:"text"`
}
