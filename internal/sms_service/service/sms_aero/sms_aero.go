package sms_aero

import (
	"fmt"
	"phone-tokens/internal/sms_service/model/interface"
	"strconv"
	"time"

	smsaero_golang "github.com/smsaero/smsaero_golang/smsaero"
)

// AeroSmsResponse
// implementation of sms response from adapter.
// /**
type AeroSmsResponse struct {
	sms    smsaero_golang.SendSms
	status smsaero_golang.HlrCheck
	err    error
}

func (r *AeroSmsResponse) GetID() string {
	return strconv.Itoa(r.sms.Id)
}

func (r *AeroSmsResponse) IsSuccess() bool {
	return r.err == nil
}

type AeroSmsElem struct {
	smsaero_golang.SendSms
	dateAnswer int
}

type AeroSmsStatus struct {
	status smsaero_golang.HlrCheck
}

func (r *AeroSmsStatus) GetStatus() string {
	return r.status.ExtendHlrStatus
}

// AeroService
// realization of sms adapter - sms aero service
// /**
type AeroService struct {
	Email  string
	ApiKey string
	Client *smsaero_golang.Client
}

func NewAeroService(email string, apiKey string) *AeroService {
	client := smsaero_golang.NewSmsAeroClient(
		email, apiKey,
		smsaero_golang.WithTimeout(time.Second*10),
		smsaero_golang.WithPhoneValidation(true),
	)
	return &AeroService{email, apiKey, client}
}

func (s *AeroService) SendSms(number int, text string) (_interface.SmsResponse, error) {
	sendResult, err := s.Client.SendSms(number, text)
	if err != nil {
		return nil, err
	}
	response := AeroSmsResponse{sms: sendResult, err: err}
	return &response, nil
}

// GetSmsList
// sms aero непонятно что возвращает в списке смсок. буквально просто мапу интерфейсов
// здесь я его привожу к мапе строк
func (s *AeroService) GetSmsList() ([]map[string]string, error) {
	sms, err := s.Client.SmsList()
	if err != nil {
		return nil, err
	}
	smsRaw := sms.(map[string]interface{})

	smsList := make([]map[string]string, 0)

	for k, v := range smsRaw {
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

func (s *AeroService) GetSmsStatus(id int) (_interface.SmsStatus, error) {
	status, err := s.Client.HlrStatus(id)
	if err != nil {
		return nil, err
	}
	var aeroResponse AeroSmsStatus
	aeroResponse.status = status

	return &aeroResponse, nil
}
