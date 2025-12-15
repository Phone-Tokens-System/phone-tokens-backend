package sms_aero

type SmsStatus int

const (
	InQueue      SmsStatus = 0
	Delivered    SmsStatus = 1
	NotDelivered SmsStatus = 2
	Passed       SmsStatus = 3
	Pending      SmsStatus = 4
	Rejected     SmsStatus = 6
	InModeration SmsStatus = 8
)
