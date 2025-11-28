package model

type SmsResponse interface {
	GetID() string
	IsSuccess() bool
}
