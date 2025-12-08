package sms_aero

type SmsStatus int

const (
	inQueue      SmsStatus = 0
	delivered    SmsStatus = 1
	notDelivered SmsStatus = 2
	passed       SmsStatus = 3
	pending      SmsStatus = 4
	rejected     SmsStatus = 6
	inModeration SmsStatus = 8
)
