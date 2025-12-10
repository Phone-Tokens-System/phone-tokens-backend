package sms_aero

import (
	"fmt"
	"phone-tokens/internal/model"
	"strconv"
	"time"

	smsaero_golang "github.com/smsaero/smsaero_golang/smsaero"
)

// AeroSmsResponse
// implementation of sms_service response from adapter.
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
// realization of sms_service adapter - sms_service aero sms_service
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

func (s *AeroService) SendSms(number int, text string) (model.SmsResponse, error) {
	sendResult, err := s.Client.SendSms(number, text)
	if err != nil {
		return model.SmsResponse{}, err
	}
	response := model.SmsResponse{
		Id:           sendResult.Id,
		From:         sendResult.From,
		Number:       sendResult.Number,
		Text:         sendResult.Text,
		Status:       sendResult.Status,
		ExtendStatus: sendResult.ExtendStatus,
		Cost:         sendResult.Cost,
		DateCreated:  sendResult.DateCreate,
		DateSent:     sendResult.DateSend,
		Raw:          sendResult,
	}
	return response, nil
}

// GetSmsList
// sms_service aero непонятно что возвращает в списке смсок. буквально просто мапу интерфейсов
// здесь я его привожу к мапе строк
func (s *AeroService) GetSmsList() ([]model.SmsResponse, error) {
	sms, err := s.Client.SmsList()
	if err != nil {
		return nil, err
	}
	smsRaw := sms.(map[string]interface{})

	//smsList := make([]map[string]string, 0)
	smsResponses := make([]model.SmsResponse, 0)
	for k, v := range smsRaw {
		if k == "links" || k == "totalCount" {
			continue
		}

		itemMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		id, err := strconv.Atoi(fmt.Sprintf("%v", itemMap["id"]))
		if err != nil {
			continue
		}
		status, err := strconv.Atoi(fmt.Sprintf("%v", itemMap["status"]))
		cost, err := strconv.ParseFloat(fmt.Sprintf("%v", itemMap["cost"]), 64)
		dateCreated, err := strconv.Atoi(fmt.Sprintf("%v", itemMap["date_created"]))
		dateSent, err := strconv.Atoi(fmt.Sprintf("%v", itemMap["date_sent"]))
		smsResponse := model.SmsResponse{
			Id:           id,
			From:         fmt.Sprintf("%v", itemMap["from"]),
			Number:       fmt.Sprintf("%v", itemMap["number"]),
			Text:         fmt.Sprintf("%v", itemMap["text"]),
			Status:       status,
			ExtendStatus: fmt.Sprintf("%v", itemMap["extend_status"]),
			Cost:         cost,
			DateCreated:  dateCreated,
			DateSent:     dateSent,
			Raw:          itemMap,
		}
		//smsItem := make(map[string]string)
		//for key, val := range itemMap {
		//	smsItem[key] = fmt.Sprintf("%v", val)
		//}
		//smsList = append(smsList, smsItem)
		smsResponses = append(smsResponses, smsResponse)
	}

	return smsResponses, nil
}

func (s *AeroService) GetSmsStatus(id int) (model.SmsStatus, error) {
	status, err := s.Client.HlrStatus(id)
	if err != nil {
		return model.SmsStatus{}, err
	}
	var respStatus model.SmsStatus
	respStatus.Status = status.HlrStatus
	respStatus.ExtendStatus = status.ExtendHlrStatus
	respStatus.Number = status.Number

	return respStatus, nil
}
