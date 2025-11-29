package service

import (
	"fmt"
	"math/big"
	"phone-tokens/internal/certificates/service"
	"phone-tokens/internal/sms_service/model"
)

type SmsService struct {
	CertificateService service.CertificateService
	SmsAdapter         SmsAdapter
}

func NewSmsService(cs service.CertificateService, adapter SmsAdapter) *SmsService {
	return &SmsService{cs, adapter}
}

func (s *SmsService) SendSms(sms model.SmsRequest) (model.SmsResponse, error) {
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
