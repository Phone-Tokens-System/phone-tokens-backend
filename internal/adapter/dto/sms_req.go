package dto

type SmsReqFromAgent struct {
	ServiceName string `json:"service_name"`
	Certificate []byte `json:"certificate"`
	ClientToken string `json:"client_number"`
	Text        string `json:"text"`
}
