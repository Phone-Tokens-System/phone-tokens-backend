package service

import (
	"fmt"
	"math/big"
	"phone-tokens/internal/certificates/service"
	"phone-tokens/internal/sms_service/model"
	"phone-tokens/internal/sms_service/model/interface"
)

// SmsService
// Сам смс сервис. Использует смс адаптер и сервис сертификатов для валидации серта
type SmsService struct {
	CertificateService service.CertificateService
	SmsAdapter         SmsAdapter
}

func NewSmsService(cs service.CertificateService, adapter SmsAdapter) *SmsService {
	return &SmsService{cs, adapter}
}

// SendSms
// Отправление смс после проверки сертификата и разрешений по токену
func (s *SmsService) SendSms(sms model.SmsRequest) (_interface.SmsResponse, error) {
	err := s.CertificateService.VerifyCertificate(sms.Certificate)
	if err != nil {
		return nil, fmt.Errorf("certificate failed verification. %w: ", err)
	}

	agentId, err := s.CertificateService.ExtractAgentId(sms.Certificate)
	if err != nil {
		return nil, fmt.Errorf("certificate failed verification. %w: ", err)
	}

	res := s.CheckPermissions(*agentId)

	if !res {
		return nil, fmt.Errorf("permission denied")
	}

	sendSms, err := s.SmsAdapter.SendSms(sms.ClientNumber, sms.Text)
	if err != nil {
		return nil, err
	}
	return sendSms, nil
}

func (s *SmsService) CheckPermissions(agentId big.Int) bool {
	return true // check permissions. maybe in another service
}

// GetSmsStatus
// Получить статус отправленного смс.
func (s *SmsService) GetSmsStatus(id int) (_interface.SmsStatus, error) {
	response, err := s.SmsAdapter.GetSmsStatus(id)
	if err != nil {
		return nil, err
	}
	return response, nil
}
