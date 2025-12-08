package _interface

/** interface for sms response from sms adapter.
 */
type SmsResponse interface {
	GetID() string
	IsSuccess() bool
}
