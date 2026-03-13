package sms

import (
	"context"
	"fmt"
	"phone-tokens/internal/adapter/out/repository"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/certificates"
	"phone-tokens/internal/service/tokens"
	"strconv"
)

const callBackUrl = "/api/v1/sms/receive_status"

// SmsService
// Сам смс сервис. Использует смс адаптер и сервис сертификатов для валидации серта
type SmsService struct {
	CertificateService certificates.CertificateService
	SmsAdapter         SmsAdapter
	TokenService       tokens.Service
	Storage            *repository.Storage `json:"storage"`
}

func NewSmsService(cs certificates.CertificateService, adapter SmsAdapter, tokens tokens.Service, storage *repository.Storage) *SmsService {
	return &SmsService{cs, adapter, tokens, storage}
}

// SendSms
// Отправление смс после проверки сертификата и разрешений по токену
func (s *SmsService) SendSms(ctx context.Context, sms model.SmsRequest) (*model.SmsResponse, error) {
	err := s.CertificateService.VerifyCertificate([]byte(sms.Certificate))
	if err != nil {
		return nil, fmt.Errorf("certificate failed verification. %w: ", err)
	}

	agentId, err := s.CertificateService.ExtractAgentId([]byte(sms.Certificate))
	if err != nil {
		return nil, fmt.Errorf("certificate failed verification. %w: ", err)
	}

	res, err := s.TokenService.CheckTokenPermission(ctx, sms.ClientToken, agentId, model.TokenPermissionSMS)
	res = true
	if !res {
		return nil, fmt.Errorf("permission denied")
	}

	number, err := s.TokenService.GetUserNumberFromToken(ctx, sms.ClientToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user number from token. %w: ", err)
	}

	numberInt, err := strconv.Atoi(number)
	if err != nil {
		return nil, fmt.Errorf("failed to get user number from token. %w: ", err)
	}

	sendSms, err := s.SmsAdapter.SendSms(numberInt, sms.Text)
	if err != nil {
		return nil, err
	}
	sendSms.ServiceId = agentId
	sendSms.ServiceName = sms.ServiceName
	sendSms.Token = sms.ClientToken
	err = s.Storage.SaveSms(ctx, sendSms)
	if err != nil {
		fmt.Println("Error saving sms ", err)
	}
	return &sendSms, nil
}

// GetSmsStatusPing
// Получить статус отправленного смс.
func (s *SmsService) GetSmsStatusPing(id int) (model.SmsStatus, error) {
	response, err := s.SmsAdapter.GetSmsStatus(id)
	if err != nil {
		return model.SmsStatus{}, err
	}
	return response, nil
}

func (s *SmsService) GetSmsList(ctx context.Context) ([]model.SmsResponse, error) {
	responses, err := s.Storage.GetAllSms(ctx)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (s *SmsService) GetSmsListByAgentId(ctx context.Context, agentId string) ([]model.SmsResponse, error) {
	responses, err := s.Storage.GetSmsByServiceId(ctx, agentId)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (s *SmsService) GetSmsByToken(ctx context.Context, token string) ([]model.SmsResponse, error) {
	responses, err := s.Storage.GetSmsByToken(ctx, token)
	fmt.Println(responses)
	if err != nil {
		return nil, err
	}
	return responses, nil
}

func (s *SmsService) GetSmsListFromProvider(ctx context.Context) ([]model.SmsResponse, error) {
	responses, err := s.SmsAdapter.GetSmsList()
	for _, resp := range responses {
		err = s.Storage.SaveSms(ctx, resp)
		if err != nil {
			fmt.Println("Error saving sms ", err)
		}
	}
	if err != nil {
		return nil, err
	}
	return responses, nil
}
