package dto

type SmsId struct {
	Id int `json:"id"`
}

type SmsStatus struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
