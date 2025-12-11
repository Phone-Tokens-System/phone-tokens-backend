package sms_service

import (
	"context"
	"fmt"
	"math/big"
	"phone-tokens/internal/model"
	"phone-tokens/internal/service/certificates"
	"phone-tokens/internal/service/tokens"
	"strconv"
)

// SmsService
// Сам смс сервис. Использует смс адаптер и сервис сертификатов для валидации серта
type SmsService struct {
	CertificateService certificates.CertificateService
	SmsAdapter         SmsAdapter
	TokenService       tokens.Service
}

func NewSmsService(cs certificates.CertificateService, adapter SmsAdapter, tokens tokens.Service) *SmsService {
	return &SmsService{cs, adapter, tokens}
}

// SendSms
// Отправление смс после проверки сертификата и разрешений по токену
func (s *SmsService) SendSms(ctx context.Context, sms model.SmsRequest) (model.SmsResponse, error) {
	err := s.CertificateService.VerifyCertificate(sms.Certificate)
	if err != nil {
		return model.SmsResponse{}, fmt.Errorf("certificate failed verification. %w: ", err)
	}

	agentId, err := s.CertificateService.ExtractAgentId(sms.Certificate)
	if err != nil {
		return model.SmsResponse{}, fmt.Errorf("certificate failed verification. %w: ", err)
	}

	res := s.CheckPermissions(*agentId)

	if !res {
		return model.SmsResponse{}, fmt.Errorf("permission denied")
	}

	number, err := s.TokenService.GetUserNumberFromToken(ctx, sms.ClientToken)
	if err != nil {
		return model.SmsResponse{}, fmt.Errorf("failed to get user number from token. %w: ", err)
	}

	numberInt, err := strconv.Atoi(number)
	if err != nil {
		return model.SmsResponse{}, fmt.Errorf("failed to get user number from token. %w: ", err)
	}

	sendSms, err := s.SmsAdapter.SendSms(numberInt, sms.Text)
	if err != nil {
		return model.SmsResponse{}, err
	}

	sendSms.ServiceName = sms.ServiceName
	return sendSms, nil
}

func (s *SmsService) CheckPermissions(agentId big.Int) bool {
	return true // check permissions. maybe in another sms_service
}

// GetSmsStatus
// Получить статус отправленного смс.
func (s *SmsService) GetSmsStatus(id int) (model.SmsStatus, error) {
	response, err := s.SmsAdapter.GetSmsStatus(id)
	if err != nil {
		return model.SmsStatus{}, err
	}
	return response, nil
}

func (s *SmsService) GetSmsList() ([]model.SmsResponse, error) {
	responses, err := s.SmsAdapter.GetSmsList()
	if err != nil {
		return nil, err
	}
	return responses, nil
}
