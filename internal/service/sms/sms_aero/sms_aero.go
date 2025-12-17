package sms_aero

import (
	"encoding/json"
	"fmt"
	"phone-tokens/internal/model"
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

type AeroSmsElem struct {
	smsaero_golang.SendSms
	dateAnswer int
}

// AeroService
// realization of sms adapter - sms aero sms
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
	rawBytes, err := json.Marshal(sendResult)
	if err != nil {
		rawBytes = nil
	}
	response := model.SmsResponse{
		ExternalId:   string(rune(sendResult.Id)),
		From:         sendResult.From,
		Text:         sendResult.Text,
		Status:       sendResult.Status,
		ExtendStatus: sendResult.ExtendStatus,
		Cost:         sendResult.Cost,
		DateCreated:  sendResult.DateCreate,
		DateSent:     sendResult.DateSend,
		Raw:          rawBytes,
	}
	return response, nil
}

// GetSmsList
// sms aero непонятно что возвращает в списке смсок. буквально просто мапу интерфейсов
// здесь я его привожу к мапе строк
func (s *AeroService) GetSmsList() ([]model.SmsResponse, error) {
	fmt.Println("GetSmsList")
	sms, err := s.Client.SmsList()
	fmt.Println(sms)
	if err != nil {
		return nil, err
	}
	smsRaw := sms.(map[string]interface{})
	fmt.Println(smsRaw)
	//smsList := make([]map[string]string, 0)
	smsResponses := make([]model.SmsResponse, 0)
	for k, v := range smsRaw {
		fmt.Println("keys")
		fmt.Println(k, v)
		if k == "links" || k == "totalCount" {
			continue
		}

		itemMap, ok := v.(map[string]interface{})
		fmt.Println("map")
		fmt.Println(itemMap)
		if !ok {
			fmt.Println("nooooooooooo")
			continue
		}
		id := fmt.Sprintf("%v", itemMap["id"])
		status, _ := strconv.Atoi(fmt.Sprintf("%v", itemMap["status"]))
		cost, _ := strconv.ParseFloat(fmt.Sprintf("%v", itemMap["cost"]), 64)
		dateCreated, _ := strconv.Atoi(fmt.Sprintf("%v", itemMap["date_created"]))
		dateSent, _ := strconv.Atoi(fmt.Sprintf("%v", itemMap["date_sent"]))

		rawMap, err := json.Marshal(itemMap)
		if err != nil {
			rawMap = nil
		}
		smsResponse := model.SmsResponse{
			ExternalId:   id,
			From:         fmt.Sprintf("%v", itemMap["from"]),
			Text:         fmt.Sprintf("%v", itemMap["text"]),
			Status:       status,
			ExtendStatus: fmt.Sprintf("%v", itemMap["extend_status"]),
			Cost:         cost,
			DateCreated:  dateCreated,
			DateSent:     dateSent,
			Raw:          rawMap,
		}
		fmt.Println(smsResponse)
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
