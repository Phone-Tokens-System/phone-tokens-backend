package _interface

/** interface for sms_service response from sms_service adapter.
 */
type SmsResponse interface {
	GetID() string
	IsSuccess() bool
}
