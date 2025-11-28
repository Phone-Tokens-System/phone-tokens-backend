package service

import (
	"fmt"
	"phone-tokens/internal/sms_service/model"
	"strconv"
	"time"

	smsaero_golang "github.com/smsaero/smsaero_golang/smsaero"
)

type AeroService struct {
	Email  string
	ApiKey string
	Client *smsaero_golang.Client
}

type AeroSmsResponse struct {
	smsaero_golang.SendSms
}

func (r *AeroSmsResponse) GetID() string {
	return strconv.Itoa(r.Id)
}

func (r *AeroSmsResponse) IsSuccess() bool {
	return r.IsSuccess()
}

func NewAeroService(email string, apiKey string) *AeroService {
	client := smsaero_golang.NewSmsAeroClient(
		email, apiKey,
		smsaero_golang.WithTimeout(time.Second*10),
		smsaero_golang.WithTest(true),
		smsaero_golang.WithPhoneValidation(true),
	)
	return &AeroService{email, apiKey, client}
}

func (s *AeroService) SendSms(number int, text string) (model.SmsResponse, error) {
	sendResult, err := s.Client.SendSms(number, text)
	if err != nil {
		return nil, err
	}
	response := AeroSmsResponse{sendResult}
	return &response, nil
}

func main() {
	email := "email"
	apiKey := "token"
	service := NewAeroService(email, apiKey)
	sms, err := service.SendSms(43534534, "hello world")
	if err != nil {
		return
	}
	fmt.Println(sms)
}
