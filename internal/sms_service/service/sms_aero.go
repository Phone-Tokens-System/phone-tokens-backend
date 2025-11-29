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
	sms    smsaero_golang.SendSms
	status smsaero_golang.HlrCheck
	err    error
}

func (r *AeroSmsResponse) GetID() string {
	return strconv.Itoa(r.sms.Id)
}

type AeroSmsElem struct {
	smsaero_golang.SendSms
	dateAnswer int
}

func (r *AeroSmsResponse) IsSuccess() bool {
	return r.err == nil
}

func NewAeroService(email string, apiKey string) *AeroService {
	client := smsaero_golang.NewSmsAeroClient(
		email, apiKey,
		smsaero_golang.WithTimeout(time.Second*10),
		smsaero_golang.WithPhoneValidation(true),
	)
	return &AeroService{email, apiKey, client}
}

func (s *AeroService) SendSms(number int, text string) (model.SmsResponse, error) {
	sendResult, err := s.Client.SendSms(number, text)
	if err != nil {
		return nil, err
	}
	status, err := s.Client.HlrStatus(sendResult.Id)
	if err != nil {
		return nil, err
	}
	response := AeroSmsResponse{sendResult, status, err}
	return &response, nil
}

func (s *AeroService) GetSmsList() ([]map[string]string, error) {
	sms, err := s.Client.SmsList()
	if err != nil {
		return nil, err
	}
	smsRaw := sms.(map[string]interface{})

	smsList := make([]map[string]string, 0)

	for k, v := range smsRaw {
		// пропускаем служебные поля вроде "links" и "totalCount"
		if k == "links" || k == "totalCount" {
			continue
		}

		itemMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		smsItem := make(map[string]string)
		for key, val := range itemMap {
			smsItem[key] = fmt.Sprintf("%v", val)
		}
		smsList = append(smsList, smsItem)
	}

	return smsList, nil
}
